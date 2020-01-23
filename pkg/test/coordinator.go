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
	"github.com/onosproject/onos-test/pkg/cluster"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/util/logging"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
	"sync"
	"time"
)

// newCoordinator returns a new test coordinator
func newCoordinator(config *Config) (*Coordinator, error) {
	kubeAPI, err := kube.GetAPI(config.ID)
	if err != nil {
		return nil, err
	}
	return &Coordinator{
		client: kubeAPI.Clientset(),
		config: config,
	}, nil
}

// Coordinator coordinates workers for suites of tests
type Coordinator struct {
	client *kubernetes.Clientset
	config *Config
}

// Run runs the tests
func (c *Coordinator) Run() error {
	for iteration := 1; iteration <= c.config.Iterations || c.config.Iterations < 0; iteration++ {
		suites := c.config.Suites
		if len(suites) == 0  || suites[0] == "" {
			// No suite specified - run all of them
			suites = make([]string, 0, len(Registry.tests))
			for suite := range Registry.tests {
				suites = append(suites, suite)
			}
		}
		workers := make([]*WorkerTask, len(suites))
		for i, suite := range suites {
			jobID := newJobID(c.config.ID+"-"+strconv.Itoa(iteration), suite)
			config := &Config{
				ID:              jobID,
				Image:           c.config.Image,
				ImagePullPolicy: c.config.ImagePullPolicy,
				Suites:          []string{suite},
				Tests:           c.config.Tests,
				Env:             c.config.Env,
				Iterations:      c.config.Iterations,
			}
			testCluster, err := cluster.NewCluster(config.ID)
			if err != nil {
				return err
			}
			worker := &WorkerTask{
				client:  c.client,
				cluster: testCluster,
				config:  config,
			}
			workers[i] = worker
		}
		err := runWorkers(workers)
		if err != nil {
			return err
		}
	}
	return nil
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
	cluster *cluster.Cluster
	config  *Config
}

// Run runs the worker job
func (t *WorkerTask) Run() (int, error) {
	// Start the job
	err := t.start()
	if err != nil {
		_ = t.tearDown()
		return 0, err
	}

	// Get the stream of logs for the pod
	pod, err := t.getPod(func(pod corev1.Pod) bool {
		return len(pod.Status.ContainerStatuses) > 0 &&
			pod.Status.ContainerStatuses[0].State.Running != nil
	})
	if err != nil {
		_ = t.tearDown()
		return 0, err
	} else if pod == nil {
		_ = t.tearDown()
		return 0, errors.New("cannot locate test pod")
	}

	req := t.client.CoreV1().Pods(t.config.ID).GetLogs(pod.Name, &corev1.PodLogOptions{
		Follow: true,
	})
	reader, err := req.Stream()
	if err != nil {
		_ = t.tearDown()
		return 0, err
	}
	defer reader.Close()

	// Stream the logs to stdout
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logging.Print(scanner.Text())
	}

	// Get the exit message and code
	_, status, err := t.getStatus()
	if err != nil {
		_ = t.tearDown()
		return 0, err
	}

	// Tear down the cluster if necessary
	_ = t.tearDown()
	return status, nil
}

// start starts the worker job
func (t *WorkerTask) start() error {
	if err := t.cluster.Create(); err != nil {
		return err
	}
	if err := t.startTest(); err != nil {
		return err
	}
	if err := t.awaitTestJobRunning(); err != nil {
		return err
	}
	return nil
}

// startTest starts running a test job
func (t *WorkerTask) startTest() error {
	if err := t.createTestJob(); err != nil {
		return err
	}
	if err := t.awaitTestJobRunning(); err != nil {
		return err
	}
	return nil
}

// createTestJob creates the job to run tests
func (t *WorkerTask) createTestJob() error {
	zero := int32(0)
	one := int32(1)

	env := t.config.ToEnv()
	env[kube.NamespaceEnv] = t.config.ID
	env[testContextEnv] = string(testContextWorker)
	env[testJobEnv] = t.config.ID

	envVars := make([]corev1.EnvVar, 0, len(env))
	for key, value := range env {
		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	batchJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      t.config.ID,
			Namespace: t.config.ID,
			Annotations: map[string]string{
				"job": t.config.ID,
			},
		},
		Spec: batchv1.JobSpec{
			Parallelism:  &one,
			Completions:  &one,
			BackoffLimit: &zero,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"job": t.config.ID,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: t.config.ID,
					RestartPolicy:      corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:            "test",
							Image:           t.config.Image,
							ImagePullPolicy: t.config.ImagePullPolicy,
							Env:             envVars,
						},
					},
				},
			},
		},
	}

	if t.config.Timeout > 0 {
		timeoutSeconds := int64(t.config.Timeout / time.Second)
		batchJob.Spec.ActiveDeadlineSeconds = &timeoutSeconds
	}
	_, err := t.client.BatchV1().Jobs(t.config.ID).Create(batchJob)
	return err
}

// awaitTestJobRunning blocks until the test job creates a pod in the RUNNING state
func (t *WorkerTask) awaitTestJobRunning() error {
	for {
		pod, err := t.getPod(func(pod corev1.Pod) bool {
			return len(pod.Status.ContainerStatuses) > 0 &&
				pod.Status.ContainerStatuses[0].Ready
		})
		if err != nil {
			return err
		} else if pod != nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// getStatus gets the status message and exit code of the given test
func (t *WorkerTask) getStatus() (string, int, error) {
	for {
		pod, err := t.getPod(func(pod corev1.Pod) bool {
			return len(pod.Status.ContainerStatuses) > 0 &&
				pod.Status.ContainerStatuses[0].State.Terminated != nil
		})
		if err != nil {
			return "", 0, err
		} else if pod != nil {
			state := pod.Status.ContainerStatuses[0].State
			if state.Terminated != nil {
				return state.Terminated.Message, int(state.Terminated.ExitCode), nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// getPod finds the Pod for the given test
func (t *WorkerTask) getPod(predicate func(pod corev1.Pod) bool) (*corev1.Pod, error) {
	pods, err := t.client.CoreV1().Pods(t.config.ID).List(metav1.ListOptions{
		LabelSelector: "job=" + t.config.ID,
	})
	if err != nil {
		return nil, err
	} else if len(pods.Items) > 0 {
		for _, pod := range pods.Items {
			if predicate(pod) {
				return &pod, nil
			}
		}
	}
	return nil, nil
}

// tearDown tears down the job
func (t *WorkerTask) tearDown() error {
	return t.cluster.Delete()
}
