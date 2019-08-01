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
		tests:      make(map[string]Test),
		TestSuites: make(map[string]TestSuite),
	}
}

//NewTestSuite returns a pointer to a new TestSuite
func NewTestSuite(name string) *TestSuite {
	return &TestSuite{
		name:  name,
		tests: make(map[string]Test),
	}
}

// Test is a test function
type Test func(t *testing.T)

// TestRegistry contains a mapping of named test groups
type TestRegistry struct {
	tests      map[string]Test
	TestSuites map[string]TestSuite
}

//TestSuite to run multiple tests
type TestSuite struct {
	name  string
	tests map[string]Test
}

//RegisterTest registers a test to the registry
func (r *TestRegistry) RegisterTest(name string, test Test, suites []*TestSuite) {
	r.tests[name] = test

	for _, suite := range suites {
		//fmt.Println("Registering test: ", name, "on suite: ",suite.name)
		suite.registerTest(name, test)
	}
}

//RegisterTest registers a test to a test group
func (r *TestSuite) registerTest(name string, test Test) {
	r.tests[name] = test
}

//RegisterTestSuite registers test suite into the registry
func (r *TestRegistry) RegisterTestSuite(testGroup TestSuite) {
	r.TestSuites[testGroup.name] = testGroup
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

// GetTestNames returns a slice of test names
func (r *TestSuite) GetTestNames() []string {
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
