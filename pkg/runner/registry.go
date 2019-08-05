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

package runner

import (
	"sort"
	"testing"
)

// NewRegistry returns a pointer to a new TestRegistry
func NewRegistry() *TestRegistry {
	return &TestRegistry{
		tests:       make(map[string]Test),
		benchmarks:  make(map[string]Benchmark),
		TestSuites:  make(map[string]TestSuite),
		BenchSuites: make(map[string]BenchSuite),
	}
}

// NewTestSuite returns a pointer to a new TestSuite
func NewTestSuite(name string) *TestSuite {
	return &TestSuite{
		name:  name,
		tests: make(map[string]Test),
	}
}

// NewBenchSuite returns a pointer to a new TestSuite
func NewBenchSuite(name string) *BenchSuite {
	return &BenchSuite{
		name:       name,
		benchmarks: make(map[string]Benchmark),
	}
}

// Test is a test function
type Test func(t *testing.T)

// Benchmark is a benchmark function
type Benchmark func(t *testing.B)

// TestRegistry contains a mapping of named test groups
type TestRegistry struct {
	tests       map[string]Test
	benchmarks  map[string]Benchmark
	TestSuites  map[string]TestSuite
	BenchSuites map[string]BenchSuite
}

// RegisterTest adds a test to the registry
func (r *TestRegistry) RegisterTest(name string, test Test, suites []*TestSuite) {
	r.tests[name] = test
	for _, suite := range suites {
		suite.registerTest(name, test)
	}
}

// RegisterBench adds a benchmark to the registry
func (r *TestRegistry) RegisterBench(name string, benchmark Benchmark, suites []*BenchSuite) {
	r.benchmarks[name] = benchmark
	for _, suite := range suites {
		suite.registerBenchmark(name, benchmark)
	}
}

// GetTestNames returns a slice of test names
func (r *TestRegistry) GetTestNames() []string {
	names := make([]string, 0, len(r.tests))
	for name := range r.tests {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}

// GetTestSuiteNames returns a slice of test names
func (r *TestRegistry) GetTestSuiteNames() []string {
	names := make([]string, 0, len(r.TestSuites))
	for name := range r.TestSuites {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}

// RegisterTestSuite registers test suite into the registry
func (r *TestRegistry) RegisterTestSuite(testGroup TestSuite) {
	r.TestSuites[testGroup.name] = testGroup
}

// GetBenchmarkNames returns a slice of benchmark names
func (r *TestRegistry) GetBenchmarkNames() []string {
	names := make([]string, 0, len(r.benchmarks))
	for name := range r.benchmarks {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}

// GetBenchSuiteNames returns a slice of test names
func (r *TestRegistry) GetBenchSuiteNames() []string {
	names := make([]string, 0, len(r.BenchSuites))
	for name := range r.BenchSuites {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}

// RegisterBenchSuite adds benchmark suite into the registry
func (r *TestRegistry) RegisterBenchSuite(suite BenchSuite) {
	r.BenchSuites[suite.name] = suite
}

// TestSuite to run multiple tests
type TestSuite struct {
	name  string
	tests map[string]Test
}

// RegisterTest registers a test to a test group
func (s *TestSuite) registerTest(name string, test Test) {
	s.tests[name] = test
}

// GetTestNames returns a slice of test names
func (s *TestSuite) GetTestNames() []string {
	names := make([]string, 0, len(s.tests))
	for name := range s.tests {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}

// BenchSuite to run multiple tests
type BenchSuite struct {
	name       string
	benchmarks map[string]Benchmark
}

// RegisterTest registers a test to a test group
func (s *BenchSuite) registerBenchmark(name string, benchmark Benchmark) {
	s.benchmarks[name] = benchmark
}

// GetBenchNames returns a slice of test names
func (s *BenchSuite) GetBenchNames() []string {
	names := make([]string, 0, len(s.benchmarks))
	for name := range s.benchmarks {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}
