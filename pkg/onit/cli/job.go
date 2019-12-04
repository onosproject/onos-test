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

package cli

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/test"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"time"
)

// cliContext indicates the context in which the CLI is running
type cliContext string

const cliContextEnv = "ONIT_CONTEXT"

const (
	localContext cliContext = "local"
	k8sContext   cliContext = "kubernetes"
)

// getContext returns the current CLI context
func getContext() cliContext {
	context := os.Getenv(cliContextEnv)
	if context == "" {
		return localContext
	}
	return cliContext(context)
}

// runInNewCluster wraps the given command to create a new cluster and execute the command inside the cluster
func runInNewCluster(command func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// If running inside k8s, simply run the command
		if getContext() == k8sContext {
			return command(cmd, args)
		}

		// Create the cluster namespace
		kubeAPI, err := kube.GetAPI(getCluster(cmd))
		if err != nil {
			return err
		}
		// Create the cluster
		c, err := test.NewCluster(kubeAPI.Namespace())
		if err != nil {
			return err
		}
		if err := c.Create(); err != nil {
			return err
		}
		command := newCommand(kubeAPI)
		return command.run(os.Args)
	}
}

// runInCluster wraps the given command to execute it inside the cluster
func runInCluster(command func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// If running inside k8s, simply run the command
		if getContext() == k8sContext {
			return command(cmd, args)
		}

		// Get the cluster API and execute the command
		kubeAPI, err := kube.GetAPI(getCluster(cmd))
		if err != nil {
			return err
		}
		command := newCommand(kubeAPI)
		return command.run(os.Args)
	}
}

// newCommandJob returns a new command job
func newCommand(api kube.API) *commandJob {
	return &commandJob{
		api:  api,
		name: fmt.Sprintf("onit-%s", random.NewPetName(2)),
	}
}

// commandJob manages the execution of commands inside the k8s cluster
type commandJob struct {
	api  kube.API
	name string
}

// run runs the command
func (j *commandJob) run(command []string) error {
	if err := j.startJob(command); err != nil {
		return err
	}
	if err := j.streamLogs(); err != nil {
		return err
	}
	code, err := j.getResult()
	if err != nil {
		return err
	} else if code > 0 {
		os.Exit(code)
	}
	return nil
}

// startJob starts the job for the given command
func (j *commandJob) startJob(command []string) error {
	var one int32 = 1
	timeout := int64((5 * time.Minute) / time.Second)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.name,
			Namespace: j.api.Namespace(),
			Labels: map[string]string{
				"type": "command",
				"job":  j.name,
			},
		},
		Spec: batchv1.JobSpec{
			Parallelism:           &one,
			Completions:           &one,
			ActiveDeadlineSeconds: &timeout,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"type": "command",
						"job":  j.name,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: j.api.Namespace(),
					RestartPolicy:      corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:            "onit",
							Image:           "onosproject/onit:latest",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         append([]string{"onit"}, command[1:]...),
							Env: []corev1.EnvVar{
								{
									Name:  cliContextEnv,
									Value: string(k8sContext),
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := j.api.Clientset().BatchV1().Jobs(j.api.Namespace()).Create(job)
	return err
}

// getPod gets the job pod
func (j *commandJob) getPod(predicate func(pod corev1.Pod) bool) (*corev1.Pod, error) {
	for {
		pods, err := j.api.Clientset().CoreV1().Pods(j.api.Namespace()).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf("type=command,job=%s", j.name),
		})
		if err != nil {
			return nil, err
		}
		if len(pods.Items) > 0 {
			pod := pods.Items[0]
			if predicate(pod) {
				return &pod, nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// streamLogs streams the job's logs to stdout
func (j *commandJob) streamLogs() error {
	pod, err := j.getPod(func(pod corev1.Pod) bool {
		return len(pod.Status.ContainerStatuses) > 0 &&
			pod.Status.ContainerStatuses[0].State.Running != nil
	})
	if err != nil {
		return err
	}

	req := j.api.Clientset().CoreV1().Pods(j.api.Namespace()).GetLogs(pod.Name, &corev1.PodLogOptions{
		Follow: true,
	})
	reader, err := req.Stream()
	if err != nil {
		return err
	}
	defer reader.Close()

	// Stream the logs to stdout
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logging.Print(scanner.Text())
	}
	return nil
}

// getResult gets the result of the job
func (j *commandJob) getResult() (int, error) {
	pod, err := j.getPod(func(pod corev1.Pod) bool {
		return len(pod.Status.ContainerStatuses) > 0 &&
			pod.Status.ContainerStatuses[0].State.Terminated != nil
	})
	if err != nil {
		return 0, err
	}

	state := pod.Status.ContainerStatuses[0].State
	if state.Terminated != nil {
		return int(state.Terminated.ExitCode), nil
	}
	return 1, errors.New("test job is not complete")
}
