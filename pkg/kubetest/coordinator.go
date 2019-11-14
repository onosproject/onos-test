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
	"bufio"
	"errors"
	"fmt"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/util/logging"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"sync"
)

// newTestCoordinator returns a new test coordinator
func newTestCoordinator(test *TestConfig) (Coordinator, error) {
	if kubeAPI, err := kube.GetAPI(test.TestID); err != nil {
		return nil, err
	} else {
		return &TestCoordinator{
			client: kubeAPI.Clientset(),
			test:   test,
		}, nil
	}
}

// newBenchmarkCoordinator returns a new benchmark coordinator
func newBenchmarkCoordinator(test *TestConfig) (Coordinator, error) {
	if kubeAPI, err := kube.GetAPI(test.TestID); err != nil {
		return nil, err
	} else {
		return &BenchmarkCoordinator{
			client: kubeAPI.Clientset(),
			test:   test,
		}, nil
	}
}

// Coordinator coordinates workers for tests and benchmarks
type Coordinator interface {
	// Run runs the coordinator
	Run() error
}

// TestCoordinator coordinates workers for suites of tests
type TestCoordinator struct {
	client *kubernetes.Clientset
	test   *TestConfig
}

// Run runs the tests
func (c *TestCoordinator) Run() error {
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
				Teardown:   c.test.Teardown,
			}
			job := &TestJob{
				cluster: &TestCluster{
					client:    c.client,
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
			Test:       c.test.Test,
			Timeout:    c.test.Timeout,
			PullPolicy: c.test.PullPolicy,
			Teardown:   c.test.Teardown,
		}
		job := &TestJob{
			cluster: &TestCluster{
				client:    c.client,
				namespace: config.TestID,
			},
			test: config,
		}
		jobs = append(jobs, job)
	}
	return runJobs(jobs, c.client)
}

// BenchmarkCoordinator coordinates workers for suites of benchmarks
type BenchmarkCoordinator struct {
	client *kubernetes.Clientset
	test   *TestConfig
}

// Run runs the tests
func (c *BenchmarkCoordinator) Run() error {
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
				Teardown:   c.test.Teardown,
			}
			job := &TestJob{
				cluster: &TestCluster{
					client:    c.client,
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
			Test:       c.test.Test,
			Timeout:    c.test.Timeout,
			PullPolicy: c.test.PullPolicy,
			Teardown:   c.test.Teardown,
		}
		job := &TestJob{
			cluster: &TestCluster{
				client:    c.client,
				namespace: config.TestID,
			},
			test: config,
		}
		jobs = append(jobs, job)
	}
	return runJobs(jobs, c.client)
}

// runJobs runs the given test jobs
func runJobs(jobs []*TestJob, client *kubernetes.Clientset) error {
	// Start jobs in separate goroutines
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	errChan := make(chan error, len(jobs))
	codeChan := make(chan int, len(jobs))
	for _, job := range jobs {
		wg.Add(1)
		go func(job *TestJob) {
			// Start the job
			err := job.start()
			if err != nil {
				errChan <- err
				_ = job.tearDown()
				return
			}

			// Get the stream of logs for the pod
			pod, err := job.getPod()
			if err != nil {
				errChan <- err
				_ = job.tearDown()
				return
			} else if pod == nil {
				errChan <- errors.New("cannot locate test pod")
				_ = job.tearDown()
				return
			}

			req := client.CoreV1().Pods(job.test.TestID).GetLogs(pod.Name, &corev1.PodLogOptions{
				Follow: true,
			})
			reader, err := req.Stream()
			if err != nil {
				errChan <- err
				if job.test.Teardown {
					_ = job.tearDown()
				}
				return
			}
			defer reader.Close()

			// Stream the logs to stdout
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				mu.Lock()
				logging.Print(scanner.Text())
				mu.Unlock()
			}

			// Get the exit message and code
			_, status, err := job.getStatus()
			if err != nil {
				errChan <- err
				if job.test.Teardown {
					_ = job.tearDown()
				}
				return
			} else {
				codeChan <- status
			}

			// Tear down the cluster if necessary
			if job.test.Teardown {
				_ = job.tearDown()
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
