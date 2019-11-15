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

package atomix

import (
	"github.com/onosproject/onos-test/pkg/onit"
	"github.com/onosproject/onos-test/pkg/onit/setup"
)

// testSuite is a suite of tests for Atomix primitives
type testSuite struct {
	onit.TestSuite
}

// SetupTestSuite sets up the Atomix test suite
func (s *testSuite) SetupTestSuite() {
	setup.Database().
		SetPartitions(3).
		SetNodesPerPartition(3)
	setup.SetupOrDie()
}

// SmokeTestSuite is a suite of tests for Atomix primitives
type SmokeTestSuite struct {
	testSuite
}

// HATestSuite is a suite of HA tests for Atomix primitives
type HATestSuite struct {
	testSuite
}

// BenchmarkSuite is a suite of benchmarks for Atomix primitives
type BenchmarkSuite struct {
	onit.BenchmarkSuite
}

// SetupBenchmarkSuite sets up the Atomix benchmark suite
func (s *BenchmarkSuite) SetupBenchmarkSuite() {
	setup.Database().
		SetPartitions(3).
		SetNodesPerPartition(3)
	setup.SetupOrDie()
}
