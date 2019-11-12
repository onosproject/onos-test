// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubetest

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/util/k8s"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

// newTestCoordinator returns a new test coordinator
func newTestCoordinator(test *TestConfig) (Coordinator, error) {
	client, err := k8s.GetClient()
	if err != nil {
		return nil, err
	}
	return &TestCoordinator{
		client: client,
		test:   test,
	}, nil
}

// newBenchmarkCoordinator returns a new benchmark coordinator
func newBenchmarkCoordinator(test *TestConfig) (Coordinator, error) {
	client, err := k8s.GetClient()
	if err != nil {
		return nil, err
	}
	return &BenchmarkCoordinator{
		client: client,
		test:   test,
	}, nil
}

// Coordinator coordinates workers for tests and benchmarks
type Coordinator interface {
	// Run runs the coordinator
	Run() error
}

// TestCoordinator coordinates workers for suites of tests
type TestCoordinator struct {
	client client.Client
	test   *TestConfig
}

// Run runs the tests
func (c *TestCoordinator) Run() error {
	client, err := k8s.GetClientset()
	if err != nil {
		return err
	}

	jobs := make([]*TestJob, 0)
	if c.test.Suite == "" {
		for suite := range Registry.tests {
			config := &TestConfig{
				TestID:     newJobID(c.test.TestID, suite),
				Type:       c.test.Type,
				Image:      c.test.Image,
				Suite:      suite,
				Timeout:    c.test.Timeout,
				PullPolicy: c.test.PullPolicy,
			}
			job := &TestJob{
				cluster: &TestCluster{
					client:    client,
					namespace: config.TestID,
				},
				test: config,
			}
			jobs = append(jobs, job)
		}
	} else {
		config := &TestConfig{
			TestID:     newJobID(c.test.TestID, c.test.Suite),
			Type:       c.test.Type,
			Image:      c.test.Image,
			Suite:      c.test.Suite,
			Timeout:    c.test.Timeout,
			PullPolicy: c.test.PullPolicy,
		}
		job := &TestJob{
			cluster: &TestCluster{
				client:    client,
				namespace: config.TestID,
			},
			test: config,
		}
		jobs = append(jobs, job)
	}
	return runJobs(jobs)
}

// BenchmarkCoordinator coordinates workers for suites of benchmarks
type BenchmarkCoordinator struct {
	client client.Client
	test   *TestConfig
}

// Run runs the tests
func (c *BenchmarkCoordinator) Run() error {
	client, err := k8s.GetClientset()
	if err != nil {
		return err
	}

	jobs := make([]*TestJob, 0)
	if c.test.Suite == "" {
		for suite := range Registry.benchmarks {
			config := &TestConfig{
				TestID:     newJobID(c.test.TestID, suite),
				Type:       c.test.Type,
				Image:      c.test.Image,
				Suite:      suite,
				Timeout:    c.test.Timeout,
				PullPolicy: c.test.PullPolicy,
			}
			job := &TestJob{
				cluster: &TestCluster{
					client:    client,
					namespace: config.TestID,
				},
				test: config,
			}
			jobs = append(jobs, job)
		}
	} else {
		config := &TestConfig{
			TestID:     newJobID(c.test.TestID, c.test.Suite),
			Type:       c.test.Type,
			Image:      c.test.Image,
			Suite:      c.test.Suite,
			Timeout:    c.test.Timeout,
			PullPolicy: c.test.PullPolicy,
		}
		job := &TestJob{
			cluster: &TestCluster{
				client:    client,
				namespace: config.TestID,
			},
			test: config,
		}
		jobs = append(jobs, job)
	}
	return runJobs(jobs)
}

// runJobs runs the given test jobs
func runJobs(jobs []*TestJob) error {
	// Start jobs in separate goroutines
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	errChan := make(chan error, len(jobs))
	codeChan := make(chan int, len(jobs))
	for _, job := range jobs {
		wg.Add(1)
		go func(job *TestJob) {
			if output, code, err := job.Run(); err != nil {
				errChan <- err
			} else {
				mu.Lock()
				_, _ = os.Stdout.WriteString(output)
				codeChan <- code
				mu.Unlock()
			}
			wg.Done()
		}(job)
	}

	// Wait for all jobs to start before proceeding
	go func() {
		wg.Wait()
		close(errChan)
		close(codeChan)
	}()

	// If any job returned an error, return it
	for err := range errChan {
		return err
	}

	// If any job returned a non-zero exit code, exit with it
	for code := range codeChan {
		if code != 0 {
			os.Exit(code)
		}
	}
	return nil
}

// newJobID returns a new unique test job ID
func newJobID(testID, suite string) string {
	return fmt.Sprintf("%s-%s", testID, suite)
}
