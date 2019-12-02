package benchmark

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/rand"
	"math"
	"sort"
	"sync"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNNOPQRSTUVWXYZ1234567890"

// New returns a new benchmark
func New() *Benchmark {
	return &Benchmark{
		clients:     1,
		parallelism: 1,
		requests:    1,
	}
}

// Handler is a handler for a benchmark
type Handler interface {
	// Run runs an iteration of the benchmark
	Run(args ...interface{}) error
}

// Benchmark is a benchmark runner
type Benchmark struct {
	clients        int
	requests       int
	parallelism    int
	handlerFactory func() Handler
	handlerArgs    []Arg
	handlers       []Handler
}

// SetClients sets the number of clients to use
func (b *Benchmark) SetClients(clients int) *Benchmark {
	b.clients = clients
	return b
}

// SetRequests sets the total number of requests to run
func (b *Benchmark) SetRequests(iterations int) *Benchmark {
	b.requests = iterations
	return b
}

// SetParallelism sets the number of parallel requests to allow per client
func (b *Benchmark) SetParallelism(parallelism int) *Benchmark {
	b.parallelism = parallelism
	return b
}

// SetHandlerFactory sets the benchmark handler factory
func (b *Benchmark) SetHandlerFactory(f func() Handler) *Benchmark {
	b.handlerFactory = f
	return b
}

// SetHandlerArgs sets the per-iteration benchmark handler arguments
func (b *Benchmark) SetHandlerArgs(args ...Arg) *Benchmark {
	b.handlerArgs = args
	return b
}

// AddHandlerArg adds a per-iteration benchmark handler argument
func (b *Benchmark) AddHandlerArg(arg Arg) *Benchmark {
	b.handlerArgs = append(b.handlerArgs, arg)
	return b
}

// Run runs the benchmark
func (b *Benchmark) Run() {
	// Prepare the benchmark arguments
	b.prepare()

	// Warm the benchmark
	b.warm()

	// Run the benchmark
	runTime, results := b.run()

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
	println(fmt.Sprintf("Operations: %d", b.requests))
	println(fmt.Sprintf("Operations/sec: %f", float64(b.requests)/(float64(runTime)/float64(time.Second))))
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
func (b *Benchmark) prepare() {
	// Prepare the handler arguments
	for _, arg := range b.handlerArgs {
		arg.Reset()
	}

	// Create the handlers
	handlers := make([]Handler, b.clients)
	for i := 0; i < b.clients; i++ {
		handlers[i] = b.handlerFactory()
	}
	b.handlers = handlers
}

// getArgs returns the arguments for the given number of requests
func (b *Benchmark) getArgs(requests int) [][]interface{} {
	args := make([][]interface{}, requests)
	for i := 0; i < requests; i++ {
		requestArgs := make([]interface{}, len(b.handlerArgs))
		for j, arg := range b.handlerArgs {
			requestArgs[j] = arg.Next()
		}
		args[i] = requestArgs
	}
	return args
}

// warm warms up the benchmark
func (b *Benchmark) warm() {
	// Create an iteration channel and wait group and create a goroutine for each handler
	wg := &sync.WaitGroup{}
	requestCh := make(chan []interface{}, len(b.handlers)*b.parallelism)
	for _, handler := range b.handlers {
		for i := 0; i < b.parallelism; i++ {
			wg.Add(1)
			go func(handler Handler) {
				for args := range requestCh {
					_ = handler.Run(args...)
				}
				wg.Done()
			}(handler)
		}
	}

	// Create the arguments for benchmark calls
	args := b.getArgs(b.requests)

	// Record the start time and write arguments to the channel
	for i := 0; i < len(args); i++ {
		requestCh <- args[i]
	}
	close(requestCh)

	// Wait for the tests to finish and close the result channel
	wg.Wait()
}

// run runs the benchmark
func (b *Benchmark) run() (time.Duration, []time.Duration) {
	// Create an iteration channel and wait group and create a goroutine for each handler
	wg := &sync.WaitGroup{}
	requestCh := make(chan []interface{}, len(b.handlers)*b.parallelism)
	resultCh := make(chan time.Duration, b.requests)
	for _, handler := range b.handlers {
		for i := 0; i < b.parallelism; i++ {
			wg.Add(1)
			go func(handler Handler) {
				for args := range requestCh {
					start := time.Now()
					_ = handler.Run(args...)
					end := time.Now()
					resultCh <- end.Sub(start)
				}
				wg.Done()
			}(handler)
		}
	}

	// Create the arguments for benchmark calls
	args := b.getArgs(b.requests)

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

type Arg interface {
	Reset()
	Next() interface{}
}

func RandomString(count int, length int) Arg {
	return &randomStringArg{
		count:  count,
		length: length,
	}
}

type randomStringArg struct {
	count  int
	length int
	args   []string
}

func (a *randomStringArg) Reset() {
	a.args = make([]string, a.count)
	for i := 0; i < a.count; i++ {
		bytes := make([]byte, a.length)
		for j := 0; j < a.length; j++ {
			bytes[j] = chars[rand.Intn(len(chars))]
		}
		a.args[i] = string(bytes)
	}
}

func (a *randomStringArg) Next() interface{} {
	return a.args[rand.Intn(len(a.args))]
}

func RandomBytes(count int, length int) Arg {
	return &randomBytesArg{
		count:  count,
		length: length,
	}
}

type randomBytesArg struct {
	count  int
	length int
	args   [][]byte
}

func (a *randomBytesArg) Reset() {
	a.args = make([][]byte, a.count)
	for i := 0; i < a.count; i++ {
		bytes := make([]byte, a.length)
		for j := 0; j < a.length; j++ {
			bytes[j] = chars[rand.Intn(len(chars))]
		}
		a.args[i] = bytes
	}
}

func (a *randomBytesArg) Next() interface{} {
	return a.args[rand.Intn(len(a.args))]
}
