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
	"context"
	"errors"
	"github.com/ghodss/yaml"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

// TestJob manages a single test job for a suite
type TestJob struct {
	client client.Client
	test   *TestConfig
}

// Start starts the test job
func (j *TestJob) Start() error {
	if err := j.ensureNamespace(); err != nil {
		return err
	}
	return j.startTests()
}

// WaitForComplete waits for the test job to finish running
func (j *TestJob) WaitForComplete() error {
	return j.awaitTestJobComplete()
}

// GetResult gets the job result
func (j *TestJob) GetResult() (string, int, error) {
	return j.getStatus()
}

// ensureNamespace sets up the test namespace
func (j *TestJob) ensureNamespace() error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: j.test.TestID,
		},
	}
	if err := j.client.Create(context.Background(), ns); err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return j.setupRBAC()
}

// setupRBAC sets up role based access controls for the cluster
func (j *TestJob) setupRBAC() error {
	if err := j.createClusterRole(); err != nil {
		return err
	}
	if err := j.createClusterRoleBinding(); err != nil {
		return err
	}
	if err := j.createServiceAccount(); err != nil {
		return err
	}
	return nil
}

// createClusterRole creates the ClusterRole required by the Atomix controller and tests if not yet created
func (j *TestJob) createClusterRole() error {
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.test.TestID,
			Namespace: j.test.TestID,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods",
					"pods/log",
					"pods/exec",
					"services",
					"endpoints",
					"persistentvolumeclaims",
					"events",
					"configmaps",
					"secrets",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"namespaces",
				},
				Verbs: []string{
					"get",
				},
			},
			{
				APIGroups: []string{
					"apps",
				},
				Resources: []string{
					"deployments",
					"daemonsets",
					"replicasets",
					"statefulsets",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"policy",
				},
				Resources: []string{
					"poddisruptionbudgets",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"k8s.atomix.io",
				},
				Resources: []string{
					"*",
				},
				Verbs: []string{
					"*",
				},
			},
		},
	}
	if err := j.client.Create(context.Background(), role); err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createClusterRoleBinding creates the ClusterRoleBinding required by the test manager
func (j *TestJob) createClusterRoleBinding() error {
	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.test.TestID,
			Namespace: j.test.TestID,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      j.test.TestID,
				Namespace: j.test.TestID,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     j.test.TestID,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	if err := j.client.Create(context.Background(), roleBinding); err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createServiceAccount creates a ServiceAccount used by the test manager
func (j *TestJob) createServiceAccount() error {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.test.TestID,
			Namespace: j.test.TestID,
		},
	}
	if err := j.client.Create(context.Background(), serviceAccount); err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// startTests starts running a test job
func (j *TestJob) startTests() error {
	if err := j.createTestConfig(); err != nil {
		return err
	}
	if err := j.createTestJob(); err != nil {
		return err
	}
	if err := j.awaitTestJobRunning(); err != nil {
		return err
	}
	return nil
}

// createTestConfig creates a ConfigMap for the test configuration
func (j *TestJob) createTestConfig() error {
	data, err := yaml.Marshal(j.test)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.test.TestID,
			Namespace: j.test.TestID,
		},
		Data: map[string]string{
			configFile: string(data),
		},
	}
	return j.client.Create(context.Background(), cm)
}

// createTestJob creates the job to run tests
func (j *TestJob) createTestJob() error {
	zero := int32(0)
	one := int32(1)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.test.TestID,
			Namespace: j.test.TestID,
			Annotations: map[string]string{
				"test-id": j.test.TestID,
				"suite":   j.test.Suite,
			},
		},
		Spec: batchv1.JobSpec{
			Parallelism:  &one,
			Completions:  &one,
			BackoffLimit: &zero,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"type": testType,
						"test": j.test.TestID,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: j.test.TestID,
					RestartPolicy:      corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:            "test",
							Image:           j.test.Image,
							ImagePullPolicy: j.test.PullPolicy,
							Env: []corev1.EnvVar{
								{
									Name:  testContextEnv,
									Value: string(TestContextWorker),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: configPath,
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: j.test.TestID,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if j.test.Timeout > 0 {
		timeoutSeconds := int64(j.test.Timeout / time.Second)
		job.Spec.ActiveDeadlineSeconds = &timeoutSeconds
	}
	return j.client.Create(context.Background(), job)
}

// awaitTestJobRunning blocks until the test job creates a pod in the RUNNING state
func (j *TestJob) awaitTestJobRunning() error {
	for {
		pod, err := j.getPod()
		if err != nil {
			return err
		} else if pod != nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// awaitTestJobComplete blocks until the test job is complete
func (j *TestJob) awaitTestJobComplete() error {
	for {
		pod, err := j.getPod()
		if err != nil {
			return err
		} else if pod == nil {
			return errors.New("cannot locate test pod")
		}
		state := pod.Status.ContainerStatuses[0].State
		if state.Terminated != nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// getStatus gets the status message and exit code of the given pod
func (j *TestJob) getStatus() (string, int, error) {
	pod, err := j.getPod()
	if err != nil {
		return "", 0, err
	} else if pod == nil {
		return "", 0, errors.New("cannot locate test pod")
	}
	state := pod.Status.ContainerStatuses[0].State
	if state.Terminated != nil {
		return state.Terminated.Message, int(state.Terminated.ExitCode), nil
	}
	return "", 0, errors.New("test job is not complete")
}

// getPod finds the Pod for the given test
func (j *TestJob) getPod() (*corev1.Pod, error) {
	pods := &corev1.PodList{}
	opts := &client.ListOptions{
		Namespace: j.test.TestID,
		LabelSelector: labels.SelectorFromSet(map[string]string{
			"test": j.test.TestID,
		}),
	}
	if err := j.client.List(context.Background(), opts, pods); err != nil {
		return nil, err
	} else if len(pods.Items) > 0 {
		for _, pod := range pods.Items {
			if pod.Status.Phase == corev1.PodRunning && len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].Ready {
				return &pod, nil
			}
		}
		for _, pod := range pods.Items {
			if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
				return &pod, nil
			}
		}
	}
	return nil, nil
}

// TearDown tears down the job
func (j *TestJob) TearDown() error {
	ns := &corev1.Namespace{}
	name := types.NamespacedName{
		Name: j.test.TestID,
	}
	if err := j.client.Get(context.Background(), name, ns); err != nil {
		return err
	}
	if err := j.client.Delete(context.Background(), ns); err != nil {
		return nil
	}
	return nil
}
