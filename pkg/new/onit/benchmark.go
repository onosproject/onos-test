package onit

import "github.com/onosproject/onos-test/pkg/new/kubetest"

// Benchmarks is the base type for ONIT benchmark suites
type Benchmarks struct {
	*kubetest.Benchmarks
}

// SetupBenchmarkSuite sets up the ONOS cluster
func (b *Benchmarks) SetupBenchmarkSuite() {
	setupONOSBenchmark(b)
}

// BenchmarkSuite is an ONIT benchmark suite
type BenchmarkSuite interface {
	kubetest.BenchmarkSuite
}

// SetupONOSBenchmarkSuite is an interface for setting up an ONOS benchmark
type SetupONOSBenchmarkSuite interface {
	SetupONOSBenchmarkSuite(setup Setup)
}

// setupONOSBenchmark sets up the ONOS cluster for the given benchmark suite
func setupONOSBenchmark(b BenchmarkSuite) {
	if setupONOS, ok := b.(SetupONOSBenchmarkSuite); ok {
		setupONOS.SetupONOSBenchmarkSuite(NewSetup(b.KubeAPI().Config()))
	}
}
