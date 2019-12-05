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
	"errors"
	"fmt"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"time"
)

// GetClusters returns a list of test clusters
func GetClusters() ([]string, error) {
	kubeAPI, err := kube.GetAPIFromEnv()
	if err != nil {
		return nil, err
	}

	namespaces, err := kubeAPI.Clientset().CoreV1().Namespaces().List(metav1.ListOptions{})
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

// NewCluster returns a new test cluster for the given Kubernetes API
func NewCluster(namespace string) (*Cluster, error) {
	kubeAPI, err := kube.GetAPIFromEnv()
	if err != nil {
		return nil, err
	}
	return &Cluster{
		client:    kubeAPI.Clientset(),
		namespace: namespace,
	}, nil
}

// Cluster manages a test cluster
type Cluster struct {
	client    *kubernetes.Clientset
	namespace string
}

// Create creates the cluster
func (c *Cluster) Create() error {
	return c.setupNamespace()
}

// Delete deletes the cluster
func (c *Cluster) Delete() error {
	return c.teardownNamespace()
}

// setupNamespace sets up the test namespace
func (c *Cluster) setupNamespace() error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.namespace,
			Labels: map[string]string{
				"test": c.namespace,
			},
		},
	}
	step := logging.NewStep(c.namespace, "Setup namespace")
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
func (c *Cluster) setupRBAC() error {
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
func (c *Cluster) createClusterRole() error {
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
func (c *Cluster) createClusterRoleBinding() error {
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
func (c *Cluster) createServiceAccount() error {
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

// Start starts running a test job
func (c *Cluster) Start(job *Job) error {
	for i := 0; i < getBenchmarkWorkers(); i++ {
		if err := c.createWorker(i, job); err != nil {
			return err
		}
	}
	if err := c.awaitRunning(job); err != nil {
		return err
	}
	return nil
}

func getWorkerName(worker int) string {
	return fmt.Sprintf("worker-%d", worker)
}

func (c *Cluster) getWorkerAddress(worker int) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:5000", getWorkerName(worker), c.namespace)
}

// CreateWorkers creates the benchmark workers
func (c *Cluster) CreateWorkers(job *Job) error {
	for i := 0; i < getBenchmarkWorkers(); i++ {
		if err := c.createWorker(i, job); err != nil {
			return err
		}
	}
	return c.awaitRunning(job)
}

// createWorker creates the given worker
func (c *Cluster) createWorker(worker int, job *Job) error {
	envVars := []corev1.EnvVar{
		{
			Name:  benchmarkContextEnv,
			Value: string(benchmarkContextWorker),
		},
		{
			Name:  benchmarkNamespaceEnv,
			Value: c.namespace,
		},
	}
	env := job.Env
	for key, value := range env {
		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: getWorkerName(worker),
			Labels: map[string]string{
				"benchmark": job.ID,
				"worker":    fmt.Sprintf("%d", worker),
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: c.namespace,
			RestartPolicy:      corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:            "benchmark",
					Image:           job.Image,
					ImagePullPolicy: job.ImagePullPolicy,
					Env:             envVars,
					Ports: []corev1.ContainerPort{
						{
							Name:          "management",
							ContainerPort: 5000,
						},
					},
					ReadinessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(5000),
							},
						},
						InitialDelaySeconds: 2,
						PeriodSeconds:       5,
					},
				},
			},
		},
	}
	if _, err := c.client.CoreV1().Pods(c.namespace).Create(pod); err != nil {
		return err
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: getWorkerName(worker),
			Labels: map[string]string{
				"benchmark": job.ID,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"benchmark": job.ID,
				"worker":    fmt.Sprintf("%d", worker),
			},
			Ports: []corev1.ServicePort{
				{
					Name: "management",
					Port: 5000,
				},
			},
		},
	}
	if _, err := c.client.CoreV1().Services(c.namespace).Create(svc); err != nil {
		return err
	}
	return nil
}

// awaitRunning blocks until the job creates a pod in the RUNNING state
func (c *Cluster) awaitRunning(job *Job) error {
	for i := 0; i < getBenchmarkWorkers(); i++ {
		if err := c.awaitWorkerRunning(i, job); err != nil {
			return err
		}
	}
	return nil
}

// awaitWorkerRunning blocks until the given worker is running
func (c *Cluster) awaitWorkerRunning(worker int, job *Job) error {
	for {
		pod, err := c.getPod(worker, job)
		if err != nil {
			return err
		} else if pod != nil && len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].Ready {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// getWorkerConns returns the worker clients for the given benchmark
func (c *Cluster) getWorkers(job *Job) ([]WorkerServiceClient, error) {
	workers := make([]WorkerServiceClient, getBenchmarkWorkers())
	for i := 0; i < getBenchmarkWorkers(); i++ {
		worker, err := grpc.Dial(c.getWorkerAddress(i), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		workers[i] = NewWorkerServiceClient(worker)
	}
	return workers, nil
}

// GetResult gets the status message and exit code of the given test
func (c *Cluster) GetResult(job *Job) (string, int, error) {
	for i := 0; i < getBenchmarkWorkers(); i++ {
		pod, err := c.getPod(i, job)
		if err != nil {
			return "", 0, err
		}
		if pod != nil {
			state := pod.Status.ContainerStatuses[0].State
			if state.Terminated != nil {
				if state.Terminated.ExitCode > 0 {
					return state.Terminated.Message, int(state.Terminated.ExitCode), nil
				}
			} else {
				return "", 0, errors.New("test job is not complete")
			}
		} else {
			return "", 0, errors.New("test job is not complete")
		}
	}
	return "", 0, nil
}

// getPod finds the Pod for the given test
func (c *Cluster) getPod(worker int, config *Job) (*corev1.Pod, error) {
	pod, err := c.client.CoreV1().Pods(c.namespace).Get(getWorkerName(worker), metav1.GetOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return nil, err
	}
	return pod, nil
}

// teardownNamespace tears down the cluster namespace
func (c *Cluster) teardownNamespace() error {
	return c.client.CoreV1().Namespaces().Delete(c.namespace, &metav1.DeleteOptions{})
}
