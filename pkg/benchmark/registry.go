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

package benchmark

// Register registers a benchmark suite
func Register(name string, suite BenchmarkingSuite) {
	Registry.register(name, suite)
}

// Registry is the global benchmark registry
var Registry = &benchmarkRegistry{
	benchmarks: make(map[string]BenchmarkingSuite),
}

// benchmarkRegistry is a registry of runnable benchmarks
type benchmarkRegistry struct {
	benchmarks map[string]BenchmarkingSuite
}

// register registers a benchmark suite
func (s *benchmarkRegistry) register(name string, suite BenchmarkingSuite) {
	s.benchmarks[name] = suite
}
