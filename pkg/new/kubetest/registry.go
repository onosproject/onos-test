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

// RegisterTests registers a test suite
func (s *TestRegistry) RegisterTests(name string, suite TestSuite) {
	s.tests[name] = suite
}

// RegisterBenchmarks registers a benchmark suite
func (s *TestRegistry) RegisterBenchmarks(name string, suite BenchmarkSuite) {
	s.benchmarks[name] = suite
}
