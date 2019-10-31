package kubetest

// RegisterTests registers a test suite
func RegisterTests(name string, suite TestSuite) {
	Registry.RegisterTests(name, suite)
}

// RegisterBenchmarks registers a benchmark suite
func RegisterBenchmarks(name string, suite BenchmarkSuite) {
	Registry.RegisterBenchmarks(name, suite)
}

// Registry is the global test registry
var Registry = &TestRegistry{
	tests:      make(map[string]TestSuite),
	benchmarks: make(map[string]BenchmarkSuite),
}

// TestRegistry is a registry of runnable tests
type TestRegistry struct {
	tests      map[string]TestSuite
	benchmarks map[string]BenchmarkSuite
}

// Register registers a test suite
func (s *TestRegistry) RegisterTests(name string, suite TestSuite) {
	s.tests[name] = suite
}

// Register registers a benchmark suite
func (s *TestRegistry) RegisterBenchmarks(name string, suite BenchmarkSuite) {
	s.benchmarks[name] = suite
}
