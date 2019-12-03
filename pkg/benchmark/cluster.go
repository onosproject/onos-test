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
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"os"
	"sync"
	"text/tabwriter"
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
func (c *Cluster) Start(config *CoordinatorConfig) error {
	if err := c.createConfig(config); err != nil {
		return err
	}

	for i := 0; i < config.Workers; i++ {
		if err := c.createWorker(i, config); err != nil {
			return err
		}
	}
	if err := c.awaitRunning(config); err != nil {
		return err
	}
	return nil
}

// createConfig creates a ConfigMap for the test configuration
func (c *Cluster) createConfig(config *CoordinatorConfig) error {
	worker := &WorkerConfig{
		JobID:       config.JobID,
		Suite:       config.Suite,
		Parallelism: config.Parallelism,
		Args:        config.Args,
	}

	data, err := yaml.Marshal(worker)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.JobID,
			Namespace: c.namespace,
		},
		Data: map[string]string{
			configFile: string(data),
		},
	}
	_, err = c.client.CoreV1().ConfigMaps(c.namespace).Create(cm)
	return err
}

func getWorkerName(worker int) string {
	return fmt.Sprintf("worker-%d", worker)
}

func (c *Cluster) getWorkerAddress(worker int) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:5000", getWorkerName(worker), c.namespace)
}

// CreateWorkers creates the benchmark workers
func (c *Cluster) CreateWorkers(config *CoordinatorConfig) error {
	if err := c.createConfig(config); err != nil {
		return err
	}
	for i := 0; i < config.Workers; i++ {
		if err := c.createWorker(i, config); err != nil {
			return err
		}
	}
	return c.awaitRunning(config)
}

// createWorker creates the given worker
func (c *Cluster) createWorker(worker int, config *CoordinatorConfig) error {
	envVars := []corev1.EnvVar{
		{
			Name:  testContextEnv,
			Value: string(testContextWorker),
		},
		{
			Name:  testNamespaceEnv,
			Value: c.namespace,
		},
	}
	env := config.Env
	if env != nil {
		for key, value := range env {
			envVars = append(envVars, corev1.EnvVar{
				Name:  key,
				Value: value,
			})
		}
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: getWorkerName(worker),
			Labels: map[string]string{
				"benchmark": config.JobID,
				"worker":    fmt.Sprintf("%d", worker),
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: c.namespace,
			RestartPolicy:      corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:            "benchmark",
					Image:           config.Image,
					ImagePullPolicy: config.PullPolicy,
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
								Name: config.JobID,
							},
						},
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
				"benchmark": config.JobID,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"benchmark": config.JobID,
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
func (c *Cluster) awaitRunning(config *CoordinatorConfig) error {
	for i := 0; i < config.Workers; i++ {
		if err := c.awaitWorkerRunning(i, config); err != nil {
			return err
		}
	}
	return nil
}

// awaitWorkerRunning blocks until the given worker is running
func (c *Cluster) awaitWorkerRunning(worker int, config *CoordinatorConfig) error {
	for {
		pod, err := c.getPod(worker, config)
		if err != nil {
			return err
		} else if pod != nil && pod.Status.ContainerStatuses[0].Ready {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// RunBenchmarks runs the given benchmarks
func (c *Cluster) RunBenchmarks(config *CoordinatorConfig) error {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	fmt.Fprintln(writer, "BENCHMARK\tREQUESTS\tDURATION\tTHROUGHPUT\tMEAN LATENCY\tMEDIAN LATENCY\t75% LATENCY\t95% LATENCY\t99% LATENCY")

	suite := Registry.benchmarks[config.Suite]
	if config.Benchmark != "" {

	} else {
		benchmarks := getBenchmarks(suite)
		for _, benchmark := range benchmarks {
			result, err := c.runBenchmark(benchmark, config)
			if err != nil {
				return err
			}
			fmt.Fprintln(writer, fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				benchmark, result.requests, result.duration, result.throughput, result.meanLatency,
				result.latencyPercentiles[.5], result.latencyPercentiles[.75],
				result.latencyPercentiles[.95], result.latencyPercentiles[.99]))
		}
	}
	return nil
}

type result struct {
	requests           int
	duration           time.Duration
	throughput         float64
	meanLatency        time.Duration
	latencyPercentiles map[float32]time.Duration
}

// runBenchmark runs the given benchmark
func (c *Cluster) runBenchmark(benchmark string, config *CoordinatorConfig) (result, error) {
	workers, err := c.getWorkers(config)
	if err != nil {
		return result{}, err
	}

	wg := &sync.WaitGroup{}
	resultCh := make(chan *Result, len(workers))
	errCh := make(chan error, len(workers))

	for _, worker := range workers {
		wg.Add(1)
		go func(worker WorkerServiceClient, requests int) {
			result, err := worker.RunBenchmark(context.Background(), &Request{
				Benchmark: benchmark,
				Requests:  uint32(requests),
			})
			if err != nil {
				errCh <- err
			} else {
				resultCh <- result
			}
			wg.Done()
		}(worker, config.Requests/len(workers))
	}

	wg.Wait()
	close(resultCh)
	close(errCh)

	for err := range errCh {
		return result{}, err
	}

	results := make([]*Result, 0, len(workers))
	for result := range resultCh {
		results = append(results, result)
	}

	var requests uint32
	var throughputSum float64
	var latencySum time.Duration
	var latency50Sum time.Duration
	var latency75Sum time.Duration
	var latency95Sum time.Duration
	var latency99Sum time.Duration
	for result := range resultCh {
		requests += result.Requests
		throughputSum += float64(result.Requests) / (float64(result.Duration) / float64(time.Second))
		latencySum += result.Latency
		latency50Sum += result.Latency
		latency75Sum += result.Latency
		latency95Sum += result.Latency
		latency99Sum += result.Latency
	}

	throughput := throughputSum / float64(len(workers))
	duration := time.Duration(throughput * float64(requests) * float64(time.Second))
	meanLatency := time.Duration(float64(latencySum) / float64(len(workers)))
	latencyPercentiles := make(map[float32]time.Duration)
	latencyPercentiles[.5] = time.Duration(float64(latency50Sum) / float64(len(workers)))
	latencyPercentiles[.75] = time.Duration(float64(latency75Sum) / float64(len(workers)))
	latencyPercentiles[.95] = time.Duration(float64(latency95Sum) / float64(len(workers)))
	latencyPercentiles[.99] = time.Duration(float64(latency99Sum) / float64(len(workers)))

	return result{
		requests:           int(requests),
		duration:           duration,
		throughput:         throughput,
		meanLatency:        meanLatency,
		latencyPercentiles: latencyPercentiles,
	}, nil
}

// getWorkerConns returns the worker clients for the given benchmark
func (c *Cluster) getWorkers(config *CoordinatorConfig) ([]WorkerServiceClient, error) {
	workers := make([]WorkerServiceClient, config.Workers)
	for i := 0; i < config.Workers; i++ {
		worker, err := grpc.Dial(c.getWorkerAddress(i), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		workers[i] = NewWorkerServiceClient(worker)
	}
	return workers, nil
}

// getLogs gets the logs from the given pod
func (c *Cluster) getLogs(pod corev1.Pod) ([]byte, error) {
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

// GetResult gets the status message and exit code of the given test
func (c *Cluster) GetResult(config *CoordinatorConfig) (string, int, error) {
	for i := 0; i < config.Workers; i++ {
		pod, err := c.getPod(i, config)
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
func (c *Cluster) getPod(worker int, config *CoordinatorConfig) (*corev1.Pod, error) {
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
