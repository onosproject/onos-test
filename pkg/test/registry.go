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

package test

// Register registers a test suite
func Register(name string, suite TestingSuite) {
	Registry.Register(name, suite)
}

// Registry is the global test registry
var Registry = &TestRegistry{
	tests: make(map[string]TestingSuite),
}

// TestRegistry is a registry of runnable tests
type TestRegistry struct {
	tests map[string]TestingSuite
}

// Register registers a test suite
func (s *TestRegistry) Register(name string, suite TestingSuite) {
	s.tests[name] = suite
}
