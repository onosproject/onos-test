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
	"errors"
	"github.com/ghodss/yaml"
	"github.com/onosproject/onos-test/pkg/new/kubetest"
	"github.com/onosproject/onos-test/pkg/new/util/logging"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

// clusterSetup handles setting up and tearing down clusters from the CLI
type clusterSetup struct {
	clusterID string
	client    *kubernetes.Clientset
}

func (c *clusterSetup) setup() error {
	return c.setupNamespace()
}

func (c *clusterSetup) teardown() error {
	return c.teardownNamespace()
}

// setupNamespace sets up the test namespace
func (c *clusterSetup) setupNamespace() error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.clusterID,
		},
	}
	step := logging.NewStep(c.clusterID, "Create worker namespace")
	step.Start()
	_, err := c.client.CoreV1().Namespaces().Create(ns)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		step.Fail(err)
		return err
	}
	step.Complete()
	return c.setupRBAC()
}

// setupRBAC sets up role based access controls for the cluster
func (c *clusterSetup) setupRBAC() error {
	step := logging.NewStep(c.clusterID, "Set up RBAC")
	step.Start()
	if err := c.createClusterRole(); err != nil {
		step.Fail(err)
		return err
	}
	if err := c.createClusterRoleBinding(); err != nil {
		step.Fail(err)
		return err
	}
	if err := c.createServiceAccount(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// createClusterRole creates the ClusterRole required by the Atomix controller and tests if not yet created
func (c *clusterSetup) createClusterRole() error {
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.clusterID,
			Namespace: c.clusterID,
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
	_, err := c.client.RbacV1().ClusterRoles().Create(role)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createClusterRoleBinding creates the ClusterRoleBinding required by the test manager
func (c *clusterSetup) createClusterRoleBinding() error {
	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.clusterID,
			Namespace: c.clusterID,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      c.clusterID,
				Namespace: c.clusterID,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     c.clusterID,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	_, err := c.client.RbacV1().ClusterRoleBindings().Create(roleBinding)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createServiceAccount creates a ServiceAccount used by the test manager
func (c *clusterSetup) createServiceAccount() error {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.clusterID,
			Namespace: c.clusterID,
		},
	}
	_, err := c.client.CoreV1().ServiceAccounts(c.clusterID).Create(serviceAccount)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// startTest starts running a test job
func (c *clusterSetup) startTest(test *kubetest.TestConfig) error {
	if err := c.createTestConfig(test); err != nil {
		return err
	}
	if err := c.createTestJob(test); err != nil {
		return err
	}
	if err := c.awaitTestJobRunning(test); err != nil {
		return err
	}
	return nil
}

// createTestConfig creates a ConfigMap for the test configuration
func (c *clusterSetup) createTestConfig(test *kubetest.TestConfig) error {
	data, err := yaml.Marshal(test)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      test.TestID,
			Namespace: c.clusterID,
		},
		Data: map[string]string{
			"config.yaml": string(data),
		},
	}
	_, err = c.client.CoreV1().ConfigMaps(c.clusterID).Create(cm)
	return err
}

// createTestJob creates the job to run tests
func (c *clusterSetup) createTestJob(test *kubetest.TestConfig) error {
	zero := int32(0)
	one := int32(1)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      test.TestID,
			Namespace: c.clusterID,
			Annotations: map[string]string{
				"test-id": test.TestID,
				"suite":   test.Suite,
			},
		},
		Spec: batchv1.JobSpec{
			Parallelism:  &one,
			Completions:  &one,
			BackoffLimit: &zero,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"type": string(test.Type),
						"test": test.TestID,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: c.clusterID,
					RestartPolicy:      corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:            "test",
							Image:           test.Image,
							ImagePullPolicy: test.PullPolicy,
							Env: []corev1.EnvVar{
								{
									Name:  "TEST_CONTEXT",
									Value: "worker",
								},
								{
									Name:  "TEST_NAMESPACE",
									Value: c.clusterID,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/config",
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
										Name: test.TestID,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if test.Timeout > 0 {
		timeoutSeconds := int64(test.Timeout / time.Second)
		job.Spec.ActiveDeadlineSeconds = &timeoutSeconds
	}
	_, err := c.client.BatchV1().Jobs(c.clusterID).Create(job)
	return err
}

// awaitTestJobRunning blocks until the test job creates a pod in the RUNNING state
func (c *clusterSetup) awaitTestJobRunning(test *kubetest.TestConfig) error {
	for {
		pod, err := c.getPod(test)
		if err != nil {
			return err
		} else if pod != nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// awaitTestJobComplete blocks until the test job is complete
func (c *clusterSetup) awaitTestJobComplete(test *kubetest.TestConfig) error {
	for {
		pod, err := c.getPod(test)
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
func (c *clusterSetup) getStatus(test *kubetest.TestConfig) (string, int, error) {
	pod, err := c.getPod(test)
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
func (c *clusterSetup) getPod(test *kubetest.TestConfig) (*corev1.Pod, error) {
	pods, err := c.client.CoreV1().Pods(c.clusterID).List(metav1.ListOptions{
		LabelSelector: "test=" + test.TestID,
	})
	if err != nil {
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

// teardownNamespace tears down the cluster namespace
func (c *clusterSetup) teardownNamespace() error {
	return c.client.CoreV1().Namespaces().Delete(c.clusterID, &metav1.DeleteOptions{})
}
