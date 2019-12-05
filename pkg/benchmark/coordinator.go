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
	"context"
	"fmt"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"k8s.io/client-go/kubernetes"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

// newCoordinator returns a new benchmark coordinator
func newCoordinator() (*Coordinator, error) {
	kubeAPI, err := kube.GetAPI(getBenchmarkNamespace())
	if err != nil {
		return nil, err
	}
	return &Coordinator{
		client: kubeAPI.Clientset(),
	}, nil
}

// Coordinator coordinates workers for suites of benchmarks
type Coordinator struct {
	client *kubernetes.Clientset
}

// Run runs the tests
func (c *Coordinator) Run() error {
	var suites []string
	suite := getBenchmarkSuite()
	if suite == "" {
		suites = make([]string, 0, len(Registry.benchmarks))
		for suite := range Registry.benchmarks {
			suites = append(suites, suite)
		}
	} else {
		suites = []string{suite}
	}

	workers := make([]*WorkerTask, len(suites))
	for i, suite := range suites {
		jobID := newJobID(getBenchmarkJob(), suite)
		env := getBenchmarkEnv()
		env[benchmarkSuiteEnv] = suite
		job := &Job{
			ID:              jobID,
			Image:           getBenchmarkImage(),
			ImagePullPolicy: getBenchmarkImagePullPolicy(),
			Command:         os.Args,
			Env:             env,
		}

		worker := &WorkerTask{
			client: c.client,
			cluster: &Cluster{
				client:    c.client,
				namespace: jobID,
			},
			job:   job,
			suite: suite,
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
	cluster *Cluster
	job     *Job
	suite   string
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
	if err := t.cluster.CreateWorkers(t.job); err != nil {
		return err
	}
	if err := t.runBenchmarks(); err != nil {
		return err
	}
	return nil
}

// runBenchmarks runs the given benchmarks
func (t *WorkerTask) runBenchmarks() error {
	results := make([]result, 0)
	benchmark := getBenchmarkName()
	if benchmark != "" {
		step := logging.NewStep(t.job.ID, "Run benchmark %s", benchmark)
		step.Start()
		result, err := t.runBenchmark(benchmark)
		if err != nil {
			step.Fail(err)
			return err
		}
		step.Complete()
		results = append(results, result)
	} else {
		suiteStep := logging.NewStep(t.job.ID, "Run benchmark suite %s", t.suite)
		suiteStep.Start()
		suite := Registry.benchmarks[t.suite]
		benchmarks := getBenchmarks(suite)
		for _, benchmark := range benchmarks {
			benchmarkSuite := logging.NewStep(t.job.ID, "Run benchmark %s", benchmark)
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
	workers, err := t.cluster.getWorkers(t.job)
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
		}(worker, getBenchmarkRequests()/len(workers))
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
