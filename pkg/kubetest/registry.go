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
func RegisterTests(name string, suite TestingSuite) {
	Registry.RegisterTests(name, suite)
}

// RegisterBenchmarks registers a benchmark suite
func RegisterBenchmarks(name string, suite BenchmarkingSuite) {
	Registry.RegisterBenchmarks(name, suite)
}

// RegisterScripts registers a script suite
func RegisterScripts(name string, suite ScriptingSuite) {
	Registry.RegisterScripts(name, suite)
}

// Registry is the global test registry
var Registry = &TestRegistry{
	tests:      make(map[string]TestingSuite),
	benchmarks: make(map[string]BenchmarkingSuite),
	scripts:    make(map[string]ScriptingSuite),
}

// TestRegistry is a registry of runnable tests
type TestRegistry struct {
	tests      map[string]TestingSuite
	benchmarks map[string]BenchmarkingSuite
	scripts    map[string]ScriptingSuite
}

// RegisterTests registers a test suite
func (s *TestRegistry) RegisterTests(name string, suite TestingSuite) {
	s.tests[name] = suite
}

// RegisterBenchmarks registers a benchmark suite
func (s *TestRegistry) RegisterBenchmarks(name string, suite BenchmarkingSuite) {
	s.benchmarks[name] = suite
}

// RegisterScript registers a script
func (s *TestRegistry) RegisterScripts(name string, suite ScriptingSuite) {
	s.scripts[name] = suite
}
