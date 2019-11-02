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
	"github.com/onosproject/onos-test/pkg/new/onit"
	"testing"
)

// TestsOne is a test
type TestsOne struct {
	*onit.Tests
}

// SetupONOSTestSuite sets up the TestOne test suite
func (s *TestsOne) SetupONOSTestSuite(setup onit.Setup) {

}

// TestFoo is an example test case
func (s *TestsOne) TestFoo(t *testing.T) {

}

// TestsTwo is a test suite
type TestsTwo struct {
	*onit.Tests
}

// TestBar is an example test case
func (s *TestsTwo) TestBar(t *testing.T) {

}

// TestsThree is an example test
type TestsThree struct {
	*onit.Tests
}

// BenchmarksOne is an example benchmark
type BenchmarksOne struct {
	*onit.Benchmarks
}
