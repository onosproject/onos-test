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

package test

import (
	"fmt"
	"math"
	"math/rand"
	"os"
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

// BenchmarkSuite is an identifier interface for benchmark suites
type BenchmarkSuite struct{}

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

// Context provides the benchmark context
type Context struct {
	config *BenchmarkConfig
}

// GetArg gets a benchmark argument
func (c *Context) GetArg(name string) *Arg {
	if value, ok := c.config.Args[name]; ok {
		return &Arg{
			value: value,
		}
	}
	return nil
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

func newBenchmark(name string, context *Context) *Benchmark {
	return &Benchmark{
		Context: context,
		Name:    name,
	}
}

// Benchmark is a benchmark runner
type Benchmark struct {
	*Context

	// Name is the name of the benchmark
	Name string
	init interface{}
}

// Init supplies a function for initializing the benchmark client
func (b *Benchmark) Init(f interface{}) {
	b.init = f
}

// Run runs the benchmark with the given parameters
func (b *Benchmark) Run(f interface{}, params ...Param) {
	// Prepare the benchmark arguments
	handler, clients := b.prepare(b.init, f, params)

	// Warm the benchmark
	b.warm(handler, clients, params)

	// Run the benchmark
	runTime, results := b.run(handler, clients, params)

	// Calculate the total latency from latency results
	var totalLatency time.Duration
	for _, result := range results {
		totalLatency += result
	}

	// Calculate latency percentiles
	meanLatency := time.Duration(int64(totalLatency) / int64(len(results)))
	latency1 := results[int(math.Max(float64(len(results)/100)-1, 0))]
	latency5 := results[int(math.Max(float64(len(results)/20)-1, 0))]
	latency25 := results[int(math.Max(float64(len(results)/4)-1, 0))]
	latency50 := results[int(math.Max(float64(len(results)/2)-1, 0))]
	latency75 := results[int(math.Max(float64(len(results)-(len(results)/4)-1), 0))]
	latency95 := results[int(math.Max(float64(len(results)-(len(results)/20)-1), 0))]
	latency99 := results[int(math.Max(float64(len(results)-(len(results)/100)-1), 0))]

	// Output the test results
	println(fmt.Sprintf("Duration: %s", runTime.String()))
	println(fmt.Sprintf("Operations: %d", b.config.Requests))
	println(fmt.Sprintf("Operations/sec: %f", float64(b.config.Requests)/(float64(runTime)/float64(time.Second))))
	println(fmt.Sprintf("Mean latency: %s", meanLatency.String()))
	println(fmt.Sprintf("1st percentile latency: %s", latency1.String()))
	println(fmt.Sprintf("5th percentile latency: %s", latency5.String()))
	println(fmt.Sprintf("25th percentile latency: %s", latency25.String()))
	println(fmt.Sprintf("50th percentile latency: %s", latency50.String()))
	println(fmt.Sprintf("75th percentile latency: %s", latency75.String()))
	println(fmt.Sprintf("90th percentile latency: %s", latency95.String()))
	println(fmt.Sprintf("99th percentile latency: %s", latency99.String()))
}

// prepare prepares the benchmark
func (b *Benchmark) prepare(init interface{}, next interface{}, params []Param) (func(...interface{}), []interface{}) {
	initV := reflect.ValueOf(init)
	initF := func() (interface{}, error) {
		vrets := initV.Call([]reflect.Value{})
		if len(vrets) != 2 {
			panic("expected two return values for Init func")
		}

		client := vrets[0].Interface()
		var err error
		if vrets[1].Interface() != nil {
			err = vrets[1].Interface().(error)
		}
		return client, err
	}

	nextV := reflect.ValueOf(next)
	nextF := func(args ...interface{}) {
		vargs := make([]reflect.Value, len(args))
		ft := nextV.Type()
		for i := 0; i < len(args); i++ {
			if args[i] != nil {
				vargs[i] = reflect.ValueOf(args[i])
			} else {
				vargs[i] = reflect.Zero(ft.In(i))
			}
		}
		_ = nextV.Call(vargs)
	}

	// Prepare the benchmark parameters
	for _, param := range params {
		param.Reset()
	}

	// Create the clients
	clients := make([]interface{}, b.config.Clients)
	for i := 0; i < b.config.Clients; i++ {
		client, err := initF()
		if err != nil {
			panic(err)
		}
		clients[i] = client
	}
	return nextF, clients
}

// getArgs returns the arguments for the given number of requests
func (b *Benchmark) getArgs(params []Param) [][]interface{} {
	args := make([][]interface{}, b.config.Requests)
	for i := 0; i < b.config.Requests; i++ {
		requestArgs := make([]interface{}, len(params))
		for j, arg := range params {
			requestArgs[j] = arg.Next()
		}
		args[i] = requestArgs
	}
	return args
}

// warm warms up the benchmark
func (b *Benchmark) warm(f func(...interface{}), clients []interface{}, params []Param) {
	// Create an iteration channel and wait group and create a goroutine for each client
	wg := &sync.WaitGroup{}
	requestCh := make(chan []interface{}, b.config.Clients*b.config.Parallelism)
	for _, client := range clients {
		for i := 0; i < b.config.Parallelism; i++ {
			wg.Add(1)
			go func(client interface{}) {
				for args := range requestCh {
					f(append([]interface{}{client}, args...)...)
				}
				wg.Done()
			}(client)
		}
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
func (b *Benchmark) run(f func(...interface{}), clients []interface{}, params []Param) (time.Duration, []time.Duration) {
	// Create an iteration channel and wait group and create a goroutine for each client
	wg := &sync.WaitGroup{}
	requestCh := make(chan []interface{}, b.config.Clients*b.config.Parallelism)
	resultCh := make(chan time.Duration, b.config.Requests)
	for _, client := range clients {
		for i := 0; i < b.config.Parallelism; i++ {
			wg.Add(1)
			go func(client interface{}) {
				for args := range requestCh {
					start := time.Now()
					f(append([]interface{}{client}, args...)...)
					end := time.Now()
					resultCh <- end.Sub(start)
				}
				wg.Done()
			}(client)
		}
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
	results := make([]time.Duration, 0, b.config.Requests)
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
	// Reset resets the benchmark parameter
	Reset()

	// Next returns the next instance of the benchmark parameter
	Next() interface{}
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

func (a *RandomStringParam) Reset() {
	a.args = make([]string, a.count)
	for i := 0; i < a.count; i++ {
		bytes := make([]byte, a.length)
		for j := 0; j < a.length; j++ {
			bytes[j] = chars[rand.Intn(len(chars))]
		}
		a.args[i] = string(bytes)
	}
}

func (a *RandomStringParam) Next() interface{} {
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

func (a *RandomBytesParam) Reset() {
	a.args = make([][]byte, a.count)
	for i := 0; i < a.count; i++ {
		bytes := make([]byte, a.length)
		for j := 0; j < a.length; j++ {
			bytes[j] = chars[rand.Intn(len(chars))]
		}
		a.args[i] = bytes
	}
}

func (a *RandomBytesParam) Next() interface{} {
	return a.args[rand.Intn(len(a.args))]
}

type internalBenchmark struct {
	name string
	f    func(*Benchmark)
}

// RunBenchmarks runs a benchmark suite
func RunBenchmarks(suite BenchmarkingSuite, config *BenchmarkConfig) {
	suiteSetupDone := false

	context := &Context{
		config: config,
	}

	methodFinder := reflect.TypeOf(suite)
	benchmarks := []internalBenchmark{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := benchmarkFilter(method.Name, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}
		if !ok {
			continue
		}
		if !suiteSetupDone {
			if setupBenchmarkSuite, ok := suite.(SetupBenchmarkSuite); ok {
				setupBenchmarkSuite.SetupBenchmarkSuite(context)
			}
			defer func() {
				if tearDownBenchmarkSuite, ok := suite.(TearDownBenchmarkSuite); ok {
					tearDownBenchmarkSuite.TearDownBenchmarkSuite(context)
				}
			}()
			suiteSetupDone = true
		}
		benchmark := internalBenchmark{
			name: method.Name,
			f: func(b *Benchmark) {
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
			},
		}
		benchmarks = append(benchmarks, benchmark)
	}
	runBenchmarks(benchmarks, context)
}

// runBenchmark runs a benchmark
func runBenchmarks(benchmarks []internalBenchmark, context *Context) {
	for _, benchmark := range benchmarks {
		println(benchmark.name)
		benchmark.f(newBenchmark(benchmark.name, context))
	}
}

// benchmarkFilter filters benchmark method names
func benchmarkFilter(name string, config *BenchmarkConfig) (bool, error) {
	if ok, _ := regexp.MatchString("^Benchmark", name); !ok {
		return false, nil
	}
	if config.Benchmark != "" {
		return config.Benchmark == name, nil
	}
	return true, nil
}
