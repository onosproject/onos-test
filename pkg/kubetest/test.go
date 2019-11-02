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

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/util/k8s"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }

// Tests is a suite of tests run on a single cluster
type Tests struct {
	*assert.Assertions
}

// Run runs the tests
func (s *Tests) Run(t *testing.T) {
	RunTests(t, s)
}

// TestSuite is an identifier interface for test suites
type TestSuite interface {
	// Run runs the test suite
	Run(t *testing.T)
}

// SetupTestSuite is an interface for setting up a suite of tests
type SetupTestSuite interface {
	SetupTestSuite(client.Client)
}

// SetupTest is an interface for setting up individual tests
type SetupTest interface {
	SetupTest(client.Client)
}

// TearDownTestSuite is an interface for tearing down a suite of tests
type TearDownTestSuite interface {
	TearDownTestSuite(client.Client)
}

// TearDownTest is an interface for tearing down individual tests
type TearDownTest interface {
	TearDownTest(client.Client)
}

// BeforeTest is an interface for executing code before every test
type BeforeTest interface {
	BeforeTest(testName string)
}

// AfterTest is an interface for executing code after every test
type AfterTest interface {
	AfterTest(testName string)
}

func failTestOnPanic(t *testing.T) {
	r := recover()
	if r != nil {
		t.Errorf("test panicked: %v\n%s", r, debug.Stack())
		t.FailNow()
	}
}

// RunTests runs a test suite
func RunTests(t *testing.T, suite TestSuite) {
	defer failTestOnPanic(t)

	client, err := k8s.GetClient()
	if err != nil {
		panic(err)
	}

	suiteSetupDone := false

	methodFinder := reflect.TypeOf(suite)
	tests := []testing.InternalTest{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := testFilter(method.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}
		if !ok {
			continue
		}
		if !suiteSetupDone {
			if setupTestSuite, ok := suite.(SetupTestSuite); ok {
				setupTestSuite.SetupTestSuite(client)
			}
			defer func() {
				if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
					tearDownTestSuite.TearDownTestSuite(client)
				}
			}()
			suiteSetupDone = true
		}
		test := testing.InternalTest{
			Name: method.Name,
			F: func(t *testing.T) {
				defer failTestOnPanic(t)

				if setupTestSuite, ok := suite.(SetupTest); ok {
					setupTestSuite.SetupTest(client)
				}
				if beforeTestSuite, ok := suite.(BeforeTest); ok {
					beforeTestSuite.BeforeTest(method.Name)
				}
				defer func() {
					if afterTestSuite, ok := suite.(AfterTest); ok {
						afterTestSuite.AfterTest(method.Name)
					}
					if tearDownTestSuite, ok := suite.(TearDownTest); ok {
						tearDownTestSuite.TearDownTest(client)
					}
				}()
				method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(t)})
			},
		}
		tests = append(tests, test)
	}
	runTests(t, tests)
}

// runTest runs a test
func runTests(b testing.TB, tests []testing.InternalTest) {
	r, ok := b.(testRunner)
	if !ok { // backwards compatibility with Go 1.6 and below
		testing.RunTests(allTestsFilter, tests)
		return
	}

	for _, test := range tests {
		r.Run(test.Name, test.F)
	}
}

// testFilter filters test method names
func testFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Test", name); !ok {
		return false, nil
	}
	return true, nil
}

// testRunner is an interface for running a test
type testRunner interface {
	Run(name string, f func(t *testing.T)) bool
}
