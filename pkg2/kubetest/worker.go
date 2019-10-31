package kubetest

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg2/util/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

// newTestWorker returns a new test worker
func newTestWorker(test *TestConfig) (*TestWorker, error) {
	client, err := k8s.GetClient()
	if err != nil {
		return nil, err
	}
	return &TestWorker{
		client: client,
		test:   test,
	}, nil
}

// TestWorker runs a test job
type TestWorker struct {
	client client.Client
	test   *TestConfig
}

// Run runs a test
func (w *TestWorker) Run() error {
	switch w.test.Type {
	case TestTypeTest:
		return w.runTest()
	case TestTypeBenchmark:
		return w.runBenchmark()
	}
	return nil
}

// runTest runs a test
func (w *TestWorker) runTest() error {
	test := Registry.tests[w.test.Suite]
	if test == nil {
		return fmt.Errorf("unknown test suite %s", w.test.Suite)
	}

	tests := []testing.InternalTest{
		{
			Name: w.test.Suite,
			F: func(t *testing.T) {
				test.Run(t)
			},
		},
	}

	testing.Main(func(_, _ string) (bool, error) { return true, nil }, tests, nil, nil)
	return nil
}

// runBenchmark runs a benchmark
func (w *TestWorker) runBenchmark() error {
	benchmark := Registry.benchmarks[w.test.Suite]
	if benchmark == nil {
		return fmt.Errorf("unknown benchmark suite %s", w.test.Suite)
	}

	benchmarks := []testing.InternalBenchmark{
		{
			Name: w.test.Suite,
			F: func(b *testing.B) {
				benchmark.Run(b)
			},
		},
	}

	testing.Main(func(_, _ string) (bool, error) { return true, nil }, nil, benchmarks, nil)
	return nil
}
