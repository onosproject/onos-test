package onit

import "github.com/onosproject/onos-test/pkg/new/kubetest"

// RegisterTests registers a test suite
func RegisterTests(name string, suite TestSuite) {
	kubetest.RegisterTests(name, suite)
}

// RegisterBenchmarks registers a benchmark suite
func RegisterBenchmarks(name string, suite BenchmarkSuite) {
	kubetest.RegisterBenchmarks(name, suite)
}
