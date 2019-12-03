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

package benchmark

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/kube"
	"k8s.io/client-go/kubernetes"
	"os"
	"sync"
)

// newCoordinator returns a new benchmark coordinator
func newCoordinator(config *CoordinatorConfig) (*Coordinator, error) {
	kubeAPI, err := kube.GetAPI(config.JobID)
	if err != nil {
		return nil, err
	}
	return &Coordinator{
		client: kubeAPI.Clientset(),
		config: config,
	}, nil
}

// Coordinator coordinates workers for suites of benchmarks
type Coordinator struct {
	client *kubernetes.Clientset
	config *CoordinatorConfig
}

// Run runs the tests
func (c *Coordinator) Run() error {
	jobs := make([]*Job, 0)
	if c.config.Suite == "" {
		for suite := range Registry.benchmarks {
			config := &CoordinatorConfig{
				JobID:       newJobID(c.config.JobID, suite),
				Image:       c.config.Image,
				Timeout:     c.config.Timeout,
				PullPolicy:  c.config.PullPolicy,
				Teardown:    c.config.Teardown,
				Suite:       suite,
				Workers:     c.config.Workers,
				Parallelism: c.config.Parallelism,
				Requests:    c.config.Requests,
				Args:        c.config.Args,
			}
			job := &Job{
				cluster: &Cluster{
					client:    c.client,
					namespace: config.JobID,
				},
				config: config,
			}
			jobs = append(jobs, job)
		}
	} else {
		config := &CoordinatorConfig{
			JobID:       newJobID(c.config.JobID, c.config.Suite),
			Image:       c.config.Image,
			Timeout:     c.config.Timeout,
			PullPolicy:  c.config.PullPolicy,
			Teardown:    c.config.Teardown,
			Suite:       c.config.Suite,
			Benchmark:   c.config.Benchmark,
			Workers:     c.config.Workers,
			Parallelism: c.config.Parallelism,
			Requests:    c.config.Requests,
			Args:        c.config.Args,
		}
		job := &Job{
			cluster: &Cluster{
				client:    c.client,
				namespace: config.JobID,
			},
			config: config,
		}
		jobs = append(jobs, job)
	}
	return runJobs(jobs)
}

// runJobs runs the given test jobs
func runJobs(jobs []*Job) error {
	// Start jobs in separate goroutines
	wg := &sync.WaitGroup{}
	errChan := make(chan error, len(jobs))
	codeChan := make(chan int, len(jobs))
	for _, job := range jobs {
		wg.Add(1)
		go func(job *Job) {
			// Start the job
			err := job.run()
			if err != nil {
				errChan <- err
				_ = job.tearDown()
				return
			}

			// Tear down the cluster if necessary
			if job.config.Teardown {
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
