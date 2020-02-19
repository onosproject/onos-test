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
	"bufio"
	"context"
	"fmt"
	"github.com/onosproject/onos-test/pkg/cluster"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/registry"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

// newCoordinator returns a new benchmark coordinator
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

// Coordinator coordinates workers for suites of benchmarks
type Coordinator struct {
	client *kubernetes.Clientset
	config *Config
}

// Run runs the tests
func (c *Coordinator) Run() error {
	var suites []string
	if c.config.Suite == "" {
		suites = registry.GetBenchmarkSuites()
	} else {
		suites = []string{c.config.Suite}
	}

	workers := make([]*WorkerTask, len(suites))
	for i, suite := range suites {
		jobID := newJobID(c.config.ID, suite)
		config := &Config{
			ID:              jobID,
			Image:           c.config.Image,
			ImagePullPolicy: c.config.ImagePullPolicy,
			Suite:           suite,
			Benchmark:       c.config.Benchmark,
			Workers:         c.config.Workers,
			Parallelism:     c.config.Parallelism,
			Requests:        c.config.Requests,
			Duration:        c.config.Duration,
			Args:            c.config.Args,
			Env:             c.config.Env,
		}
		benchCluster, err := cluster.NewCluster(jobID)
		if err != nil {
			return err
		}

		worker := &WorkerTask{
			client:  c.client,
			cluster: benchCluster,
			config:  config,
		}
		workers[i] = worker
	}
	return runWorkers(workers)
}

// runWorkers runs the given test jobs
func runWorkers(tasks []*WorkerTask) error {
	// Start jobs in separate goroutines
	wg := &sync.WaitGroup{}
	errChan := make(chan error, len(tasks))
	codeChan := make(chan int, len(tasks))
	for _, task := range tasks {
		wg.Add(1)
		go func(task *WorkerTask) {
			status, err := task.Run()
			if err != nil {
				errChan <- err
			} else {
				codeChan <- status
			}
			wg.Done()
		}(task)
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
	workers []WorkerServiceClient
}

// Run runs the worker job
func (t *WorkerTask) Run() (int, error) {
	// Start the job
	err := t.run()
	if err != nil {
		_ = t.tearDown()
		return 0, err
	}

	// Tear down the cluster if necessary
	_ = t.tearDown()
	return 0, nil
}

// start starts the test job
func (t *WorkerTask) run() error {
	if err := t.cluster.Create(); err != nil {
		return err
	}
	if err := t.createWorkers(); err != nil {
		return err
	}
	if err := t.runBenchmarks(); err != nil {
		return err
	}
	return nil
}

func getWorkerName(worker int) string {
	return fmt.Sprintf("worker-%d", worker)
}

func (t *WorkerTask) getWorkerAddress(worker int) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:5000", getWorkerName(worker), t.config.ID)
}

// createWorkers creates the benchmark workers
func (t *WorkerTask) createWorkers() error {
	for i := 0; i < t.config.Workers; i++ {
		if err := t.createWorker(i); err != nil {
			return err
		}
	}
	return t.awaitRunning()
}

// createWorker creates the given worker
func (t *WorkerTask) createWorker(worker int) error {
	env := t.config.ToEnv()
	env[kube.NamespaceEnv] = t.config.ID
	env[benchmarkContextEnv] = string(benchmarkContextWorker)
	env[benchmarkWorkerEnv] = fmt.Sprintf("%d", worker)
	env[benchmarkJobEnv] = t.config.ID

	envVars := make([]corev1.EnvVar, 0, len(env))
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
				"benchmark": t.config.ID,
				"worker":    fmt.Sprintf("%d", worker),
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: t.config.ID,
			RestartPolicy:      corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:            "benchmark",
					Image:           t.config.Image,
					ImagePullPolicy: t.config.ImagePullPolicy,
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
	_, err := t.client.CoreV1().Pods(t.config.ID).Create(pod)
	if err != nil {
		return err
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: getWorkerName(worker),
			Labels: map[string]string{
				"benchmark": t.config.ID,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"benchmark": t.config.ID,
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
	if _, err := t.client.CoreV1().Services(t.config.ID).Create(svc); err != nil {
		return err
	}

	go t.streamWorkerLogs(worker)
	return nil
}

// streamWorkerLogs streams the logs from the given worker
func (t *WorkerTask) streamWorkerLogs(worker int) {
	for {
		pod, err := t.getPod(worker)
		if err != nil || pod == nil {
			return
		}

		if len(pod.Status.ContainerStatuses) > 0 &&
			(pod.Status.ContainerStatuses[0].State.Running != nil ||
				pod.Status.ContainerStatuses[0].State.Terminated != nil) {
			req := t.client.CoreV1().Pods(t.config.ID).GetLogs(getWorkerName(worker), &corev1.PodLogOptions{
				Follow: true,
			})
			reader, err := req.Stream()
			if err != nil {
				return
			}
			defer reader.Close()

			// Stream the logs to stdout
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				logging.Print(scanner.Text())
			}
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// awaitRunning blocks until the job creates a pod in the RUNNING state
func (t *WorkerTask) awaitRunning() error {
	for i := 0; i < t.config.Workers; i++ {
		if err := t.awaitWorkerRunning(i); err != nil {
			return err
		}
	}
	return nil
}

// awaitWorkerRunning blocks until the given worker is running
func (t *WorkerTask) awaitWorkerRunning(worker int) error {
	for {
		pod, err := t.getPod(worker)
		if err != nil {
			return err
		} else if pod != nil && len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].Ready {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// getWorkerConns returns the worker clients for the given benchmark
func (t *WorkerTask) getWorkers() ([]WorkerServiceClient, error) {
	if t.workers != nil {
		return t.workers, nil
	}

	workers := make([]WorkerServiceClient, t.config.Workers)
	for i := 0; i < t.config.Workers; i++ {
		worker, err := grpc.Dial(t.getWorkerAddress(i), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		workers[i] = NewWorkerServiceClient(worker)
	}
	t.workers = workers
	return workers, nil
}

// getPod finds the Pod for the given test
func (t *WorkerTask) getPod(worker int) (*corev1.Pod, error) {
	pod, err := t.client.CoreV1().Pods(t.config.ID).Get(getWorkerName(worker), metav1.GetOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return nil, err
	}
	return pod, nil
}

// setupSuite sets up the benchmark suite
func (t *WorkerTask) setupSuite() error {
	workers, err := t.getWorkers()
	if err != nil {
		return err
	}

	worker := workers[0]
	_, err = worker.SetupSuite(context.Background(), &SuiteRequest{
		Suite: t.config.Suite,
		Args:  t.config.Args,
	})
	return err
}

// setupWorkers sets up the benchmark workers
func (t *WorkerTask) setupWorkers() error {
	workers, err := t.getWorkers()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	for _, worker := range workers {
		wg.Add(1)
		go func(worker WorkerServiceClient) {
			_, err = worker.SetupWorker(context.Background(), &SuiteRequest{
				Suite: t.config.Suite,
				Args:  t.config.Args,
			})
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}(worker)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}

// setupBenchmark sets up the given benchmark
func (t *WorkerTask) setupBenchmark(benchmark string) error {
	workers, err := t.getWorkers()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	for _, worker := range workers {
		wg.Add(1)
		go func(worker WorkerServiceClient) {
			_, err = worker.SetupBenchmark(context.Background(), &BenchmarkRequest{
				Suite:     t.config.Suite,
				Benchmark: benchmark,
				Args:      t.config.Args,
			})
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}(worker)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}

// runBenchmarks runs the given benchmarks
func (t *WorkerTask) runBenchmarks() error {
	// Setup the benchmark suite on one of the workers
	if err := t.setupSuite(); err != nil {
		return err
	}

	// Setup the workers
	if err := t.setupWorkers(); err != nil {
		return err
	}

	// Run the benchmarks
	results := make([]result, 0)
	if t.config.Benchmark != "" {
		step := logging.NewStep(t.config.ID, "Run benchmark %s", t.config.Benchmark)
		step.Start()
		result, err := t.runBenchmark(t.config.Benchmark)
		if err != nil {
			step.Fail(err)
			return err
		}
		step.Complete()
		results = append(results, result)
	} else {
		suiteStep := logging.NewStep(t.config.ID, "Run benchmark suite %s", t.config.Suite)
		suiteStep.Start()
		suite := registry.GetBenchmarkSuite(t.config.Suite)
		benchmarks := getBenchmarks(suite)
		for _, benchmark := range benchmarks {
			benchmarkSuite := logging.NewStep(t.config.ID, "Run benchmark %s", benchmark)
			benchmarkSuite.Start()
			result, err := t.runBenchmark(benchmark)
			if err != nil {
				benchmarkSuite.Fail(err)
				suiteStep.Fail(err)
				return err
			}
			benchmarkSuite.Complete()
			results = append(results, result)
		}
		suiteStep.Complete()
	}

	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	fmt.Fprintln(writer, "BENCHMARK\tREQUESTS\tDURATION\tTHROUGHPUT\tMEAN LATENCY\tMEDIAN LATENCY\t75% LATENCY\t95% LATENCY\t99% LATENCY")
	for _, result := range results {
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%d\t%s\t%f/sec\t%s\t%s\t%s\t%s\t%s",
			result.benchmark, result.requests, result.duration, result.throughput, result.meanLatency,
			result.latencyPercentiles[.5], result.latencyPercentiles[.75],
			result.latencyPercentiles[.95], result.latencyPercentiles[.99]))
	}
	writer.Flush()
	return nil
}

// runBenchmark runs the given benchmark
func (t *WorkerTask) runBenchmark(benchmark string) (result, error) {
	// Setup the benchmark
	if err := t.setupBenchmark(benchmark); err != nil {
		return result{}, err
	}

	workers, err := t.getWorkers()
	if err != nil {
		return result{}, err
	}

	wg := &sync.WaitGroup{}
	resultCh := make(chan *RunResponse, len(workers))
	errCh := make(chan error, len(workers))

	for _, worker := range workers {
		wg.Add(1)
		go func(worker WorkerServiceClient, requests int, duration *time.Duration) {
			result, err := worker.RunBenchmark(context.Background(), &RunRequest{
				Suite:       t.config.Suite,
				Benchmark:   benchmark,
				Requests:    uint32(requests),
				Duration:    duration,
				Parallelism: uint32(t.config.Parallelism),
				Args:        t.config.Args,
			})
			if err != nil {
				errCh <- err
			} else {
				resultCh <- result
			}
			wg.Done()
		}(worker, t.config.Requests/len(workers), t.config.Duration)
	}

	wg.Wait()
	close(resultCh)
	close(errCh)

	for err := range errCh {
		return result{}, err
	}

	var duration time.Duration
	var requests uint32
	var latencySum time.Duration
	var latency50Sum time.Duration
	var latency75Sum time.Duration
	var latency95Sum time.Duration
	var latency99Sum time.Duration
	for result := range resultCh {
		requests += result.Requests
		duration = time.Duration(math.Max(float64(duration), float64(result.Duration)))
		latencySum += result.Latency
		latency50Sum += result.Latency50
		latency75Sum += result.Latency75
		latency95Sum += result.Latency95
		latency99Sum += result.Latency99
	}

	throughput := float64(requests) / (float64(duration) / float64(time.Second))
	meanLatency := time.Duration(float64(latencySum) / float64(len(workers)))
	latencyPercentiles := make(map[float32]time.Duration)
	latencyPercentiles[.5] = time.Duration(float64(latency50Sum) / float64(len(workers)))
	latencyPercentiles[.75] = time.Duration(float64(latency75Sum) / float64(len(workers)))
	latencyPercentiles[.95] = time.Duration(float64(latency95Sum) / float64(len(workers)))
	latencyPercentiles[.99] = time.Duration(float64(latency99Sum) / float64(len(workers)))

	return result{
		benchmark:          benchmark,
		requests:           int(requests),
		duration:           duration,
		throughput:         throughput,
		meanLatency:        meanLatency,
		latencyPercentiles: latencyPercentiles,
	}, nil
}

type result struct {
	benchmark          string
	requests           int
	duration           time.Duration
	throughput         float64
	meanLatency        time.Duration
	latencyPercentiles map[float32]time.Duration
}

// tearDown tears down the job
func (t *WorkerTask) tearDown() error {
	return t.cluster.Delete()
}
