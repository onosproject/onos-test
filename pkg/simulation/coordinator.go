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

package simulation

import (
	"bufio"
	"context"
	"fmt"
	"github.com/onosproject/onos-test/pkg/cluster"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/model"
	"github.com/onosproject/onos-test/pkg/registry"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"google.golang.org/grpc"
	"io"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"os"
	"reflect"
	"sync"
	"time"
)

// newCoordinator returns a new simulation coordinator
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

// Coordinator coordinates workers for suites of simulators
type Coordinator struct {
	client *kubernetes.Clientset
	config *Config
}

// Run runs the simulations
func (c *Coordinator) Run() error {
	var suites []string
	if c.config.Simulation == "" {
		suites = make([]string, 0, len(registry.GetSimulationSuites()))
		for _, suite := range registry.GetSimulationSuites() {
			suites = append(suites, suite)
		}
	} else {
		suites = []string{c.config.Simulation}
	}

	workers := make([]*WorkerTask, len(suites))
	for i, suite := range suites {
		jobID := newJobID(c.config.ID, suite)
		config := &Config{
			ID:              jobID,
			Image:           c.config.Image,
			ImagePullPolicy: c.config.ImagePullPolicy,
			Simulation:      suite,
			Model:           c.config.Model,
			Simulators:      c.config.Simulators,
			Rate:            c.config.Rate,
			Jitter:          c.config.Jitter,
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
	workers []SimulatorServiceClient
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
	if err := t.setupSimulation(); err != nil {
		return err
	}
	if err := t.setupSimulators(); err != nil {
		return err
	}
	traces, err := t.runSimulation()
	if err != nil {
		return err
	}
	if err := t.checkModel(traces); err != nil {
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

// createWorkers creates the simulation workers
func (t *WorkerTask) createWorkers() error {
	for i := 0; i < t.config.Simulators; i++ {
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
	env[simulationContextEnv] = string(simulationContextWorker)
	env[simulationWorkerEnv] = fmt.Sprintf("%d", worker)
	env[simulationJobEnv] = t.config.ID

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
				"simulation": t.config.ID,
				"worker":     fmt.Sprintf("%d", worker),
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: t.config.ID,
			RestartPolicy:      corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:            "simulation",
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
				"simulation": t.config.ID,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"simulation": t.config.ID,
				"worker":     fmt.Sprintf("%d", worker),
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
	for i := 0; i < t.config.Simulators; i++ {
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

// getSimulators returns the worker clients for the given simulation
func (t *WorkerTask) getSimulators() ([]SimulatorServiceClient, error) {
	if t.workers != nil {
		return t.workers, nil
	}

	workers := make([]SimulatorServiceClient, t.config.Simulators)
	for i := 0; i < t.config.Simulators; i++ {
		worker, err := grpc.Dial(t.getWorkerAddress(i), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		workers[i] = NewSimulatorServiceClient(worker)
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

// setupSimulation sets up the simulation
func (t *WorkerTask) setupSimulation() error {
	workers, err := t.getSimulators()
	if err != nil {
		return err
	}

	worker := workers[0]
	_, err = worker.SetupSimulation(context.Background(), &SimulationLifecycleRequest{
		Simulation: t.config.Simulation,
		Args:       t.config.Args,
	})
	return err
}

// setupSimulators sets up the simulators
func (t *WorkerTask) setupSimulators() error {
	simulators, err := t.getSimulators()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	for i, simulator := range simulators {
		wg.Add(1)
		go func(simulator int, client SimulatorServiceClient) {
			if err := t.setupSimulator(simulator, client); err != nil {
				errCh <- err
			}
			wg.Done()
		}(i, simulator)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}

// setupSimulator sets up the given simulator
func (t *WorkerTask) setupSimulator(simulator int, client SimulatorServiceClient) error {
	step := logging.NewStep(t.config.ID, "Setup simulator %s/%d", t.config.Simulation, simulator)
	step.Start()
	request := &SimulationLifecycleRequest{
		Simulation: t.config.Simulation,
		Args:       t.config.Args,
	}
	_, err := client.SetupSimulator(context.Background(), request)
	if err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// runSimulation runs the given simulations
func (t *WorkerTask) runSimulation() ([]*model.Trace, error) {
	// Run the simulation for the configured duration
	step := logging.NewStep(t.config.ID, "Run simulation %s", t.config.Simulation)
	step.Start()

	wg := &sync.WaitGroup{}
	tracesCh := make(chan *model.Trace)
	errCh := make(chan error)
	for i := 0; i < t.config.Parallelism; i++ {
		wg.Add(1)
		go func() {
			if err := t.runSimulators(tracesCh); err != nil {
				errCh <- err
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(tracesCh)
		close(errCh)
	}()

	traces := make([]*model.Trace, 0)
	for trace := range tracesCh {
		traces = append(traces, trace)
	}

	for err := range errCh {
		step.Fail(err)
		return nil, err
	}
	step.Complete()
	return traces, nil
}

// runSimulators runs the simulation for a goroutine
func (t *WorkerTask) runSimulators(ch chan<- *model.Trace) error {
	duration := time.After(t.config.Duration)
	for {
		select {
		case <-duration:
			return nil
		case <-waitJitter(t.config.Rate, t.config.Jitter):
			t.runSimulator(t.chooseSimulator(), ch)
		}
	}
}

// runSimulator runs a random simulator
func (t *WorkerTask) runSimulator(simulator int, ch chan<- *model.Trace) {
	step := logging.NewStep(t.config.ID, "Run simulator %s/%d", t.config.Simulation, simulator)
	step.Start()

	simulators, err := t.getSimulators()
	if err != nil {
		step.Fail(err)
		return
	}

	client := simulators[simulator]
	request := &SimulateRequest{
		Simulation: t.config.Simulation,
		Method:     t.chooseMethod(),
	}
	stream, err := client.Simulate(context.Background(), request)
	if err != nil {
		step.Fail(err)
		return
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			step.Complete()
			return
		} else if err != nil {
			step.Fail(err)
			return
		}
		ch <- response.Trace
	}
}

// chooseSimulator chooses a random simulator
func (t *WorkerTask) chooseSimulator() int {
	return rand.Intn(t.config.Simulators)
}

// chooseMethod chooses a random simulator method
func (t *WorkerTask) chooseMethod() string {
	suite := registry.GetSimulationSuite(t.config.Simulation)
	methods := getMethods(suite)
	return methods[rand.Intn(len(methods))]
}

// getMethods returns a list of simulators in the given suite
func getMethods(suite SimulatingSuite) []string {
	methodFinder := reflect.TypeOf(suite)
	simulators := []string{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := simulatorFilter(method.Name)
		if ok {
			simulators = append(simulators, method.Name)
		} else if err != nil {
			panic(err)
		}
	}
	return simulators
}

// checkModel checks the given traces against the model
func (t *WorkerTask) checkModel(traces []*model.Trace) error {
	if t.config.Model == "" {
		return nil
	}

	step := logging.NewStep(t.config.ID, "Check model %s", t.config.Model)
	step.Start()

	checker, err := model.NewChecker()
	if err != nil {
		step.Fail(err)
		return err
	}

	model := model.NewModel(t.config.Model, traces)
	if err := checker.Check(model); err != nil {
		step.Fail(err)
		return err
	}
	return nil
}

// tearDown tears down the job
func (t *WorkerTask) tearDown() error {
	return t.cluster.Delete()
}

// waitJitter returns a channel that closes after time.Duration between duration and duration + maxFactor *
// duration.
func waitJitter(duration time.Duration, maxFactor float64) <-chan time.Time {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	delay := duration + time.Duration(rand.Float64()*maxFactor*float64(duration))
	return time.After(delay)
}
