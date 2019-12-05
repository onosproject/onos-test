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

package test

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"sync"
)

// newCoordinator returns a new test coordinator
func newCoordinator() (*Coordinator, error) {
	kubeAPI, err := kube.GetAPI(getTestNamespace())
	if err != nil {
		return nil, err
	}
	return &Coordinator{
		client: kubeAPI.Clientset(),
	}, nil
}

// Coordinator coordinates workers for suites of tests
type Coordinator struct {
	client *kubernetes.Clientset
}

// Run runs the tests
func (c *Coordinator) Run() error {
	var suites []string
	suite := getTestSuite()
	if suite == "" {
		suites = make([]string, 0, len(Registry.tests))
		for suite := range Registry.tests {
			suites = append(suites, suite)
		}
	} else {
		suites = []string{suite}
	}

	workers := make([]*WorkerTask, len(suites))
	for i, suite := range suites {
		jobID := newJobID(getTestJob(), suite)
		env := getTestEnv()
		env[testSuiteEnv] = suite
		job := &Job{
			ID:              jobID,
			Image:           getTestImage(),
			ImagePullPolicy: getTestImagePullPolicy(),
			Command:         os.Args,
			Env:             env,
		}
		worker := &WorkerTask{
			client: c.client,
			cluster: &Cluster{
				client:    c.client,
				namespace: jobID,
			},
			job: job,
		}
		workers[i] = worker
	}
	return runWorkers(workers)
}

// runWorkers runs the given test workers
func runWorkers(tasks []*WorkerTask) error {
	// Start jobs in separate goroutines
	wg := &sync.WaitGroup{}
	errChan := make(chan error, len(tasks))
	codeChan := make(chan int, len(tasks))
	for _, job := range tasks {
		wg.Add(1)
		go func(task *WorkerTask) {
			status, err := task.Run()
			if err != nil {
				errChan <- err
			} else {
				codeChan <- status
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

// WorkerTask manages a single test job for a test worker
type WorkerTask struct {
	client  *kubernetes.Clientset
	cluster *Cluster
	job     *Job
}

// Run runs the worker job
func (j *WorkerTask) Run() (int, error) {
	// Start the job
	err := j.start()
	if err != nil {
		_ = j.tearDown()
		return 0, err
	}

	// Get the stream of logs for the pod
	pod, err := j.getPod()
	if err != nil {
		_ = j.tearDown()
		return 0, err
	} else if pod == nil {
		_ = j.tearDown()
		return 0, errors.New("cannot locate test pod")
	}

	req := j.client.CoreV1().Pods(j.job.ID).GetLogs(pod.Name, &corev1.PodLogOptions{
		Follow: true,
	})
	reader, err := req.Stream()
	if err != nil {
		_ = j.tearDown()
		return 0, err
	}
	defer reader.Close()

	// Stream the logs to stdout
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logging.Print(scanner.Text())
	}

	// Get the exit message and code
	_, status, err := j.getStatus()
	if err != nil {
		_ = j.tearDown()
		return 0, err
	}

	// Tear down the cluster if necessary
	_ = j.tearDown()
	return status, nil
}

// start starts the worker job
func (j *WorkerTask) start() error {
	if err := j.cluster.Create(); err != nil {
		return err
	}
	if err := j.cluster.startTest(j.job); err != nil {
		return err
	}
	if err := j.cluster.awaitTestJobRunning(j.job); err != nil {
		return err
	}
	return nil
}

// getStatus gets the status message and exit code of the given pod
func (j *WorkerTask) getStatus() (string, int, error) {
	return j.cluster.getTestResult(j.job)
}

// getPod finds the Pod for the given test
func (j *WorkerTask) getPod() (*v1.Pod, error) {
	return j.cluster.getPod(j.job)
}

// tearDown tears down the job
func (j *WorkerTask) tearDown() error {
	return j.cluster.Delete()
}
