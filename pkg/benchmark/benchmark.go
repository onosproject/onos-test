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
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"
)

const warmUpDuration = 30 * time.Second
const aggBatchSize = 100

// BenchmarkingSuite is a suite of benchmarks
type BenchmarkingSuite interface{}

// Suite is an identifier interface for benchmark suites
type Suite struct{}

// SetupSuite is an interface for setting up a suite of benchmarks
type SetupSuite interface {
	SetupSuite(c *Context) error
}

// TearDownSuite is an interface for tearing down a suite of benchmarks
type TearDownSuite interface {
	TearDownSuite(c *Context) error
}

// SetupWorker is an interface for setting up individual benchmarks
type SetupWorker interface {
	SetupWorker(c *Context) error
}

// TearDownWorker is an interface for tearing down individual benchmarks
type TearDownWorker interface {
	TearDownWorker(c *Context) error
}

// SetupBenchmark is an interface for executing code before every benchmark
type SetupBenchmark interface {
	SetupBenchmark(c *Context) error
}

// TearDownBenchmark is an interface for executing code after every benchmark
type TearDownBenchmark interface {
	TearDownBenchmark(c *Context) error
}

// newContext returns a new benchmark context
func newContext(name string, args map[string]string) *Context {
	return &Context{
		Name: name,
		args: args,
	}
}

// Context provides the benchmark context
type Context struct {
	Name string
	args map[string]string
}

// GetArg gets a benchmark argument
func (c *Context) GetArg(name string) *Arg {
	if value, ok := c.args[name]; ok {
		return &Arg{
			value: value,
		}
	}
	return &Arg{}
}

// Arg is a benchmark argument
type Arg struct {
	value string
}

// Int returns the argument as an int
func (a *Arg) Int(def int) int {
	if a.value == "" {
		return def
	}
	i, err := strconv.Atoi(a.value)
	if err != nil {
		panic(err)
	}
	return i
}

// String returns the argument as a string
func (a *Arg) String(def string) string {
	if a.value == "" {
		return def
	}
	return a.value
}

// newBenchmark creates a new benchmark
func newBenchmark(requests int, duration *time.Duration, parallelism int, maxLatency *time.Duration, context *Context) *Benchmark {
	return &Benchmark{
		Context:     context,
		requests:    requests,
		duration:    duration,
		maxLatency:  maxLatency,
		parallelism: parallelism,
	}
}

// Benchmark is a benchmark runner
type Benchmark struct {
	*Context

	requests    int
	duration    *time.Duration
	parallelism int
	maxLatency  *time.Duration
}

// Run runs the benchmark with the given parameters
func (b *Benchmark) run(suite BenchmarkingSuite) (*RunResponse, error) {
	var f func() error
	methods := reflect.TypeOf(suite)
	if method, ok := methods.MethodByName(b.Name); ok {
		f = func() error {
			values := method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(b)})
			if len(values) == 0 {
				return nil
			} else if values[0].Interface() == nil {
				return nil
			}
			return values[0].Interface().(error)
		}
	} else {
		return nil, fmt.Errorf("unknown benchmark method %s", b.Name)
	}

	// Warm the benchmark
	b.warmRequests(f)

	// Run the benchmark
	requests, runTime, results := b.runRequests(f)

	// Calculate the total latency from latency results
	var totalLatency time.Duration
	for _, result := range results {
		totalLatency += result
	}

	// Calculate latency percentiles
	meanLatency := time.Duration(int64(totalLatency) / int64(len(results)))
	latency50 := results[int(math.Max(float64(len(results)/2)-1, 0))]
	latency75 := results[int(math.Max(float64(len(results)-(len(results)/4)-1), 0))]
	latency95 := results[int(math.Max(float64(len(results)-(len(results)/20)-1), 0))]
	latency99 := results[int(math.Max(float64(len(results)-(len(results)/100)-1), 0))]

	return &RunResponse{
		Requests:  uint32(requests),
		Duration:  runTime,
		Latency:   meanLatency,
		Latency50: latency50,
		Latency75: latency75,
		Latency95: latency95,
		Latency99: latency99,
	}, nil
}

// warm warms up the benchmark
func (b *Benchmark) warmRequests(f func() error) {
	// Create an iteration channel and wait group and create a goroutine for each client
	wg := &sync.WaitGroup{}
	requestCh := make(chan struct{}, b.parallelism)
	for i := 0; i < b.parallelism; i++ {
		wg.Add(1)
		go func() {
			for range requestCh {
				_ = f()
			}
			wg.Done()
		}()
	}

	// Run for the warm up duration to prepare the benchmark
	start := time.Now()
	for time.Since(start) < warmUpDuration {
		requestCh <- struct{}{}
	}
	close(requestCh)

	// Wait for the tests to finish and close the result channel
	wg.Wait()
}

// run runs the benchmark
func (b *Benchmark) runRequests(f func() error) (int, time.Duration, []time.Duration) {
	// Create an iteration channel and wait group and create a goroutine for each client
	wg := &sync.WaitGroup{}
	requestCh := make(chan struct{}, b.parallelism)
	resultCh := make(chan time.Duration, aggBatchSize)
	for i := 0; i < b.parallelism; i++ {
		wg.Add(1)
		go func() {
			for range requestCh {
				start := time.Now()
				_ = f()
				end := time.Now()
				resultCh <- end.Sub(start)
			}
			wg.Done()
		}()
	}

	// Start an aggregator goroutine
	results := make([]time.Duration, 0, aggBatchSize*aggBatchSize)
	aggWg := &sync.WaitGroup{}
	aggWg.Add(1)
	go func() {
		var total time.Duration
		var count = 0
		// Iterate through results and aggregate durations
		for duration := range resultCh {
			total += duration
			count++
			// Average out the durations in batches
			if count == aggBatchSize {
				results = append(results, total/time.Duration(count))

				// If the total number of batches reaches the batch size ^ 2, aggregate the aggregated results
				if len(results) == aggBatchSize*aggBatchSize {
					newResults := make([]time.Duration, 0, aggBatchSize*aggBatchSize)
					for _, result := range results {
						total += result
						count++
						if count == aggBatchSize {
							newResults = append(newResults, total/time.Duration(count))
							total = 0
							count = 0
						}
					}
					results = newResults
				}
				total = 0
				count = 0
			}
		}
		if count > 0 {
			results = append(results, total/time.Duration(count))
		}
		aggWg.Done()
	}()

	// Record the start time and write arguments to the channel
	start := time.Now()

	// Iterate through the request count or until the time duration has been met
	requests := 0
	for (b.requests == 0 || requests < b.requests) && (b.duration == nil || time.Since(start) < *b.duration) {
		requestCh <- struct{}{}
		requests++
	}
	close(requestCh)

	// Wait for the tests to finish and close the result channel
	wg.Wait()

	// Record the end time
	end := time.Now()
	duration := end.Sub(start)

	// Close the output channel
	close(resultCh)

	// Wait for the results to be aggregated
	aggWg.Wait()

	// Sort the aggregated results
	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})
	return requests, duration, results
}

// getBenchmarks returns a list of benchmarks in the given suite
func getBenchmarks(suite BenchmarkingSuite) []string {
	methodFinder := reflect.TypeOf(suite)
	benchmarks := []string{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := benchmarkFilter(method.Name)
		if ok {
			benchmarks = append(benchmarks, method.Name)
		} else if err != nil {
			panic(err)
		}
	}
	return benchmarks
}
