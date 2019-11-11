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
	"bytes"
	"errors"
	"github.com/ghodss/yaml"
	"github.com/onosproject/onos-test/pkg/new/util/k8s"
	"github.com/onosproject/onos-test/pkg/new/util/logging"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

// GetTestClusters returns a list of test clusters
func GetTestClusters() ([]string, error) {
	client, err := k8s.GetClientset()
	if err != nil {
		return nil, err
	}

	namespaces, err := client.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	clusters := make([]string, 0)
	for _, namespace := range namespaces.Items {
		if namespace.Labels["test"] != "" {
			clusters = append(clusters, namespace.Name)
		}
	}
	return clusters, nil
}

// NewTestCluster returns a new test cluster for the given Kubernetes API
func NewTestCluster(namespace string) *TestCluster {
	client, err := k8s.GetClientset()
	if err != nil {
		panic(err)
	}
	return &TestCluster{
		client:    client,
		namespace: namespace,
	}
}

// TestCluster manages a test cluster
type TestCluster struct {
	client    *kubernetes.Clientset
	namespace string
}

// Create creates the cluster
func (c *TestCluster) Create() error {
	return c.setupNamespace()
}

// Delete deletes the cluster
func (c *TestCluster) Delete() error {
	return c.teardownNamespace()
}

// setupNamespace sets up the test namespace
func (c *TestCluster) setupNamespace() error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.namespace,
			Labels: map[string]string{
				"test": c.namespace,
			},
		},
	}
	step := logging.NewStep(c.namespace, "Create namespace")
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
func (c *TestCluster) setupRBAC() error {
	step := logging.NewStep(c.namespace, "Set up RBAC")
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
func (c *TestCluster) createClusterRole() error {
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.namespace,
			Namespace: c.namespace,
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
					"serviceaccounts",
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
					"batch",
				},
				Resources: []string{
					"jobs",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"rbac.authorization.k8s.io",
				},
				Resources: []string{
					"clusterroles",
					"clusterrolebindings",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"apiextensions.k8s.io",
				},
				Resources: []string{
					"customresourcedefinitions",
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
func (c *TestCluster) createClusterRoleBinding() error {
	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.namespace,
			Namespace: c.namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      c.namespace,
				Namespace: c.namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     c.namespace,
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
func (c *TestCluster) createServiceAccount() error {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.namespace,
			Namespace: c.namespace,
		},
	}
	_, err := c.client.CoreV1().ServiceAccounts(c.namespace).Create(serviceAccount)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// StartTest starts running a test job
func (c *TestCluster) StartTest(test *TestConfig) error {
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
func (c *TestCluster) createTestConfig(test *TestConfig) error {
	data, err := yaml.Marshal(test)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      test.TestID,
			Namespace: c.namespace,
		},
		Data: map[string]string{
			configFile: string(data),
		},
	}
	_, err = c.client.CoreV1().ConfigMaps(c.namespace).Create(cm)
	return err
}

// createTestJob creates the job to run tests
func (c *TestCluster) createTestJob(test *TestConfig) error {
	zero := int32(0)
	one := int32(1)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      test.TestID,
			Namespace: c.namespace,
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
					ServiceAccountName: c.namespace,
					RestartPolicy:      corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:            "test",
							Image:           test.Image,
							ImagePullPolicy: test.PullPolicy,
							Env: []corev1.EnvVar{
								{
									Name:  testContextEnv,
									Value: string(testContextWorker),
								},
								{
									Name:  testNamespaceEnv,
									Value: c.namespace,
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
	_, err := c.client.BatchV1().Jobs(c.namespace).Create(job)
	return err
}

// awaitTestJobRunning blocks until the test job creates a pod in the RUNNING state
func (c *TestCluster) awaitTestJobRunning(test *TestConfig) error {
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

// AwaitTestComplete blocks until the test job is complete
func (c *TestCluster) AwaitTestComplete(test *TestConfig) error {
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

// GetTestOutput gets the output from the given test
func (c *TestCluster) GetTestOutput(test *TestConfig) ([]byte, error) {
	pod, err := c.getPod(test)
	if err != nil {
		return nil, err
	}
	return c.getLogs(*pod)
}

// getLogs gets the logs from the given pod
func (c *TestCluster) getLogs(pod corev1.Pod) ([]byte, error) {
	req := c.client.CoreV1().Pods(c.namespace).GetLogs(pod.Name, &corev1.PodLogOptions{})
	readCloser, err := req.Stream()
	if err != nil {
		return nil, err
	}

	defer readCloser.Close()

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(readCloser); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetTestResult gets the status message and exit code of the given test
func (c *TestCluster) GetTestResult(test *TestConfig) (string, int, error) {
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
func (c *TestCluster) getPod(test *TestConfig) (*corev1.Pod, error) {
	pods, err := c.client.CoreV1().Pods(c.namespace).List(metav1.ListOptions{
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
func (c *TestCluster) teardownNamespace() error {
	return c.client.CoreV1().Namespaces().Delete(c.namespace, &metav1.DeleteOptions{})
}
