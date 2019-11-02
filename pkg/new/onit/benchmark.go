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
	"github.com/onosproject/onos-test/pkg/new/kubetest"
	"github.com/onosproject/onos-test/pkg/new/onit/setup"
)

// Benchmarks is the base type for ONIT benchmark suites
type Benchmarks struct {
	*kubetest.Benchmarks
}

// SetupBenchmarkSuite sets up the ONOS cluster
func (b *Benchmarks) SetupBenchmarkSuite() {
	setupONOSBenchmark(b)
}

// BenchmarkSuite is an ONIT benchmark suite
type BenchmarkSuite interface {
	kubetest.BenchmarkSuite
}

// SetupONOSBenchmarkSuite is an interface for setting up an ONOS benchmark
type SetupONOSBenchmarkSuite interface {
	SetupONOSBenchmarkSuite(setup setup.TestSetup)
}

// setupONOSBenchmark sets up the ONOS cluster for the given benchmark suite
func setupONOSBenchmark(b BenchmarkSuite) {
	if setupONOS, ok := b.(SetupONOSBenchmarkSuite); ok {
		setupONOS.SetupONOSBenchmarkSuite(setup.New(b.KubeAPI()))
	}
}
