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
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
)

var (
	// AtomixTests is the complete Atomix test suite
	AtomixTests = runner.NewTestSuite("atomix")
	// AtomixBenchmarks benchmark suite for Atomix
	AtomixBenchmarks = runner.NewBenchSuite("atomix")
)

func init() {
	test.Registry.RegisterTestSuite(*AtomixTests)
	test.Registry.RegisterBenchSuite(*AtomixBenchmarks)
}
