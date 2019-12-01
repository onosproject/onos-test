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

package onit

import (
	"github.com/onosproject/onos-test/pkg/kubetest"
)

// TestSuite is the base type for ONIT test suites
type TestSuite struct {
	kubetest.TestSuite
}

// BenchmarkSuite is the base type for ONIT benchmark suites
type BenchmarkSuite struct {
	kubetest.BenchmarkSuite
}

// ScriptSuite is the base type for ONIT script suites
type ScriptSuite struct {
	kubetest.ScriptSuite
}
