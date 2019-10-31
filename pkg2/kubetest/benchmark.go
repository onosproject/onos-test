package kubetest

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg2/util/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

var allBenchmarksFilter = func(_, _ string) (bool, error) { return true, nil }

// Benchmarks is a suite of benchmarks run on a single cluster
type Benchmarks struct {
	*assert.Assertions
	require *require.Assertions
}

// Run runs the benchmarks
func (s *Benchmarks) Run(b *testing.B) {
	RunBenchmarks(b, s)
}

// BenchmarkSuite is an identifier interface for benchmark suites
type BenchmarkSuite interface {
	// Run runs the benchmark suite
	Run(b *testing.B)
}

// SetupBenchmarkSuite is an interface for setting up a suite of benchmarks
type SetupBenchmarkSuite interface {
	SetupBenchmarkSuite(client.Client)
}

// SetupBenchmark is an interface for setting up individual benchmarks
type SetupBenchmark interface {
	SetupBenchmark(client.Client)
}

// TearDownBenchmarkSuite is an interface for tearing down a suite of benchmarks
type TearDownBenchmarkSuite interface {
	TearDownBenchmarkSuite(client.Client)
}

// TearDownBenchmark is an interface for tearing down individual benchmarks
type TearDownBenchmark interface {
	TearDownBenchmark(client.Client)
}

// BeforeBenchmark is an interface for executing code before every benchmark
type BeforeBenchmark interface {
	BeforeBenchmark(testName string)
}

// AfterBenchmark is an interface for executing code after every benchmark
type AfterBenchmark interface {
	AfterBenchmark(testName string)
}

func failBenchmarkOnPanic(b *testing.B) {
	r := recover()
	if r != nil {
		b.Errorf("test panicked: %v\n%s", r, debug.Stack())
		b.FailNow()
	}
}

// RunBenchmarks runs a benchmark suite
func RunBenchmarks(b *testing.B, suite BenchmarkSuite) {
	defer failBenchmarkOnPanic(b)

	client, err := k8s.GetClient()
	if err != nil {
		panic(err)
	}

	suiteSetupDone := false

	methodFinder := reflect.TypeOf(suite)
	benchmarks := []testing.InternalBenchmark{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := benchmarkFilter(method.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}
		if !ok {
			continue
		}
		if !suiteSetupDone {
			if setupBenchmarkSuite, ok := suite.(SetupBenchmarkSuite); ok {
				setupBenchmarkSuite.SetupBenchmarkSuite(client)
			}
			defer func() {
				if tearDownBenchmarkSuite, ok := suite.(TearDownBenchmarkSuite); ok {
					tearDownBenchmarkSuite.TearDownBenchmarkSuite(client)
				}
			}()
			suiteSetupDone = true
		}
		benchmark := testing.InternalBenchmark{
			Name: method.Name,
			F: func(b *testing.B) {
				defer failBenchmarkOnPanic(b)

				if setupBenchmarkSuite, ok := suite.(SetupBenchmark); ok {
					setupBenchmarkSuite.SetupBenchmark(client)
				}
				if beforeBenchmarkSuite, ok := suite.(BeforeBenchmark); ok {
					beforeBenchmarkSuite.BeforeBenchmark(method.Name)
				}
				defer func() {
					if afterBenchmarkSuite, ok := suite.(AfterBenchmark); ok {
						afterBenchmarkSuite.AfterBenchmark(method.Name)
					}
					if tearDownBenchmarkSuite, ok := suite.(TearDownBenchmark); ok {
						tearDownBenchmarkSuite.TearDownBenchmark(client)
					}
				}()
				method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(b)})
			},
		}
		benchmarks = append(benchmarks, benchmark)
	}
	runBenchmarks(b, benchmarks)
}

// runBenchmark runs a benchmark
func runBenchmarks(b testing.TB, benchmarks []testing.InternalBenchmark) {
	r, ok := b.(benchmark)
	if !ok { // backwards compatibility with Go 1.6 and below
		testing.RunBenchmarks(allBenchmarksFilter, benchmarks)
		return
	}

	for _, benchmark := range benchmarks {
		r.Run(benchmark.Name, benchmark.F)
	}
}

// benchmarkFilter filters benchmark method names
func benchmarkFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Benchmark", name); !ok {
		return false, nil
	}
	return true, nil
}

// benchmark is an interface for running a benchmark
type benchmark interface {
	Run(name string, f func(b *testing.B)) bool
}
