package benchmark

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/rand"
	"sort"
	"sync"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNNOPQRSTUVWXYZ1234567890"

// New returns a new benchmark
func New() *Benchmark {
	return &Benchmark{
		parallelism: 1,
		iterations:  1,
	}
}

// Handler is a handler for a benchmark
type Handler interface {
	// Run runs an iteration of the benchmark
	Run(args ...interface{}) error
}

// Benchmark is a benchmark runner
type Benchmark struct {
	parallelism    int
	iterations     int
	handlerFactory func() Handler
	handlerArgs    []Arg
}

// SetParallelism sets the number of parallel handlers to use
func (b *Benchmark) SetParallelism(parallelism int) *Benchmark {
	b.parallelism = parallelism
	return b
}

// SetIterations sets the total number of iterations to run
func (b *Benchmark) SetIterations(iterations int) *Benchmark {
	b.iterations = iterations
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
	// Prepare the handler arguments
	for _, arg := range b.handlerArgs {
		arg.Reset()
	}

	// Create the handlers
	handlers := make([]Handler, b.parallelism)
	for i := 0; i < b.parallelism; i++ {
		handlers[i] = b.handlerFactory()
	}

	// Create an iteration channel and wait group and create a goroutine for each handler
	wg := &sync.WaitGroup{}
	itCh := make(chan []interface{}, len(handlers))
	outCh := make(chan time.Duration, b.iterations)
	for _, handler := range handlers {
		wg.Add(1)
		go func(handler Handler) {
			for args := range itCh {
				start := time.Now()
				_ = handler.Run(args...)
				end := time.Now()
				outCh <- end.Sub(start)
			}
			wg.Done()
		}(handler)
	}

	// Create the arguments for benchmark calls
	args := make([][]interface{}, b.iterations)
	for i := 0; i < b.iterations; i++ {
		itArgs := make([]interface{}, len(b.handlerArgs))
		for j, arg := range b.handlerArgs {
			itArgs[j] = arg.Next()
		}
		args[i] = itArgs
	}

	// Record the start time and write arguments to the channel
	start := time.Now()
	for i := 0; i < len(args); i++ {
		itCh <- args[i]
	}
	close(itCh)

	// Wait for the tests to finish and close the result channel
	wg.Wait()

	// Record the end time
	end := time.Now()
	duration := end.Sub(start)

	close(outCh)

	// Aggregate the results
	var totalLatency time.Duration
	results := make([]time.Duration, 0, b.iterations)
	for d := range outCh {
		totalLatency += d
		results = append(results, d)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})

	meanLatency := time.Duration(int64(totalLatency) / int64(len(results)))
	latency1 := results[(len(results) / 100) - 1]
	latency5 := results[(len(results) / 20) - 1]
	latency25 := results[(len(results) / 4) - 1]
	latency50 := results[(len(results)/2) - 1]
	latency75 := results[len(results) - (len(results)/4) - 1]
	latency95 := results[len(results) - (len(results) / 20) - 1]
	latency99 := results[len(results) - (len(results) / 100) - 1]

	// Output the test results
	println(fmt.Sprintf("Duration: %s", duration.String()))
	println(fmt.Sprintf("Operations: %d", b.iterations))
	println(fmt.Sprintf("Operations/sec: %f", float64(b.iterations)/(float64(duration)/float64(time.Second))))
	println(fmt.Sprintf("Mean latency: %s", meanLatency.String()))
	println(fmt.Sprintf("1st percentile latency: %s", latency1.String()))
	println(fmt.Sprintf("5th percentile latency: %s", latency5.String()))
	println(fmt.Sprintf("25th percentile latency: %s", latency25.String()))
	println(fmt.Sprintf("50th percentile latency: %s", latency50.String()))
	println(fmt.Sprintf("75th percentile latency: %s", latency75.String()))
	println(fmt.Sprintf("90th percentile latency: %s", latency95.String()))
	println(fmt.Sprintf("99th percentile latency: %s", latency99.String()))
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
