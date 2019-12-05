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
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNNOPQRSTUVWXYZ1234567890"

// BenchmarkingSuite is a suite of benchmarks
type BenchmarkingSuite interface{}

// Suite is an identifier interface for benchmark suites
type Suite struct{}

// SetupBenchmarkSuite is an interface for setting up a suite of benchmarks
type SetupBenchmarkSuite interface {
	SetupBenchmarkSuite(c *Context)
}

// SetupBenchmark is an interface for setting up individual benchmarks
type SetupBenchmark interface {
	SetupBenchmark(b *Benchmark)
}

// TearDownBenchmarkSuite is an interface for tearing down a suite of benchmarks
type TearDownBenchmarkSuite interface {
	TearDownBenchmarkSuite(c *Context)
}

// TearDownBenchmark is an interface for tearing down individual benchmarks
type TearDownBenchmark interface {
	TearDownBenchmark(b *Benchmark)
}

// BeforeBenchmark is an interface for executing code before every benchmark
type BeforeBenchmark interface {
	BeforeBenchmark(b *Benchmark)
}

// AfterBenchmark is an interface for executing code after every benchmark
type AfterBenchmark interface {
	AfterBenchmark(b *Benchmark)
}

// newContext returns a new benchmark context
func newContext(config *Config) *Context {
	return &Context{
		config: config,
	}
}

// Context provides the benchmark context
type Context struct {
	config *Config
}

// GetArg gets a benchmark argument
func (c *Context) GetArg(name string) *Arg {
	if value, ok := c.config.Args[name]; ok {
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

func newBenchmark(name string, requests int, context *Context) *Benchmark {
	return &Benchmark{
		Context:  context,
		requests: requests,
		Name:     name,
	}
}

// Benchmark is a benchmark runner
type Benchmark struct {
	*Context

	// Name is the name of the benchmark
	Name     string
	requests int
	result   *Result
}

// getResult returns the benchmark result
func (b *Benchmark) getResult() *Result {
	return b.result
}

// Run runs the benchmark with the given parameters
func (b *Benchmark) Run(f interface{}, params ...Param) {
	// Prepare the benchmark arguments
	handler := b.prepare(f, params)

	// Warm the benchmark
	b.warm(handler, params)

	// Run the benchmark
	runTime, results := b.run(handler, params)

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

	b.result = &Result{
		Requests:  uint32(b.requests),
		Duration:  runTime,
		Latency:   meanLatency,
		Latency50: latency50,
		Latency75: latency75,
		Latency95: latency95,
		Latency99: latency99,
	}
}

// prepare prepares the benchmark
func (b *Benchmark) prepare(next interface{}, params []Param) func(...interface{}) {
	v := reflect.ValueOf(next)
	f := func(args ...interface{}) {
		vargs := make([]reflect.Value, len(args))
		ft := v.Type()
		for i := 0; i < len(args); i++ {
			if args[i] != nil {
				vargs[i] = reflect.ValueOf(args[i])
			} else {
				vargs[i] = reflect.Zero(ft.In(i))
			}
		}
		_ = v.Call(vargs)
	}

	// Prepare the benchmark parameters
	for _, param := range params {
		param.reset()
	}
	return f
}

// getArgs returns the arguments for the given number of requests
func (b *Benchmark) getArgs(params []Param) [][]interface{} {
	args := make([][]interface{}, b.requests)
	for i := 0; i < b.requests; i++ {
		requestArgs := make([]interface{}, len(params))
		for j, arg := range params {
			requestArgs[j] = arg.next()
		}
		args[i] = requestArgs
	}
	return args
}

// warm warms up the benchmark
func (b *Benchmark) warm(f func(...interface{}), params []Param) {
	// Create an iteration channel and wait group and create a goroutine for each client
	wg := &sync.WaitGroup{}
	requestCh := make(chan []interface{}, b.config.Parallelism)
	for i := 0; i < b.config.Parallelism; i++ {
		wg.Add(1)
		go func() {
			for args := range requestCh {
				f(args...)
			}
			wg.Done()
		}()
	}

	// Create the arguments for benchmark calls
	args := b.getArgs(params)

	// Record the start time and write arguments to the channel
	for i := 0; i < len(args); i++ {
		requestCh <- args[i]
	}
	close(requestCh)

	// Wait for the tests to finish and close the result channel
	wg.Wait()
}

// run runs the benchmark
func (b *Benchmark) run(f func(...interface{}), params []Param) (time.Duration, []time.Duration) {
	// Create an iteration channel and wait group and create a goroutine for each client
	wg := &sync.WaitGroup{}
	requestCh := make(chan []interface{}, b.config.Parallelism)
	resultCh := make(chan time.Duration, b.requests)
	for i := 0; i < b.config.Parallelism; i++ {
		wg.Add(1)
		go func() {
			for args := range requestCh {
				start := time.Now()
				f(args...)
				end := time.Now()
				resultCh <- end.Sub(start)
			}
			wg.Done()
		}()
	}

	// Create the arguments for benchmark calls
	args := b.getArgs(params)

	// Record the start time and write arguments to the channel
	start := time.Now()
	for i := 0; i < len(args); i++ {
		requestCh <- args[i]
	}
	close(requestCh)

	// Wait for the tests to finish and close the result channel
	wg.Wait()

	// Record the end time
	end := time.Now()
	duration := end.Sub(start)

	// Close the output channel
	close(resultCh)

	// Aggregate the results
	results := make([]time.Duration, 0, b.requests)
	for d := range resultCh {
		results = append(results, d)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})
	return duration, results
}

// Param is an interface for benchmark parameters
type Param interface {
	// reset resets the benchmark parameter
	reset()

	// next returns the next instance of the benchmark parameter
	next() interface{}
}

// RandomString returns a random string parameter
func RandomString(count int, length int) Param {
	return &RandomStringParam{
		count:  count,
		length: length,
	}
}

// RandomStringParam is a random string parameter
type RandomStringParam struct {
	count  int
	length int
	args   []string
}

func (a *RandomStringParam) reset() {
	a.args = make([]string, a.count)
	for i := 0; i < a.count; i++ {
		bytes := make([]byte, a.length)
		for j := 0; j < a.length; j++ {
			bytes[j] = chars[rand.Intn(len(chars))]
		}
		a.args[i] = string(bytes)
	}
}

func (a *RandomStringParam) next() interface{} {
	return a.args[rand.Intn(len(a.args))]
}

// RandomBytes returns a random bytes parameter
func RandomBytes(count int, length int) Param {
	return &RandomBytesParam{
		count:  count,
		length: length,
	}
}

// RandomBytesParam is a random string parameter
type RandomBytesParam struct {
	count  int
	length int
	args   [][]byte
}

func (a *RandomBytesParam) reset() {
	a.args = make([][]byte, a.count)
	for i := 0; i < a.count; i++ {
		bytes := make([]byte, a.length)
		for j := 0; j < a.length; j++ {
			bytes[j] = chars[rand.Intn(len(chars))]
		}
		a.args[i] = bytes
	}
}

func (a *RandomBytesParam) next() interface{} {
	return a.args[rand.Intn(len(a.args))]
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

// setupSuite sets up the given benchmark suite
func setupSuite(suite BenchmarkingSuite, config *Config) {
	context := newContext(config)
	if setupBenchmarkSuite, ok := suite.(SetupBenchmarkSuite); ok {
		setupBenchmarkSuite.SetupBenchmarkSuite(context)
	}
}

// runBenchmark runs a benchmark method
func runBenchmark(benchmark string, requests int, suite BenchmarkingSuite, config *Config) (*Result, error) {
	methods := reflect.TypeOf(suite)
	method, ok := methods.MethodByName(benchmark)
	if !ok {
		return nil, fmt.Errorf("cannot find benchmark method %s", benchmark)
	}

	println(benchmark)
	context := newContext(config)
	b := newBenchmark(benchmark, requests, context)
	if setupBenchmarkSuite, ok := suite.(SetupBenchmark); ok {
		setupBenchmarkSuite.SetupBenchmark(b)
	}
	if beforeBenchmarkSuite, ok := suite.(BeforeBenchmark); ok {
		beforeBenchmarkSuite.BeforeBenchmark(b)
	}
	defer func() {
		if afterBenchmarkSuite, ok := suite.(AfterBenchmark); ok {
			afterBenchmarkSuite.AfterBenchmark(b)
		}
		if tearDownBenchmarkSuite, ok := suite.(TearDownBenchmark); ok {
			tearDownBenchmarkSuite.TearDownBenchmark(b)
		}
	}()
	method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(b)})
	return b.getResult(), nil
}

// benchmarkFilter filters benchmark method names
func benchmarkFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Benchmark", name); !ok {
		return false, nil
	}
	return true, nil
}
