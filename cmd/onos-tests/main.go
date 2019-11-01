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

package main

import (
	"github.com/onosproject/onos-test/pkg/kubetest"
	_ "github.com/onosproject/onos-test/test/api"
	_ "github.com/onosproject/onos-test/test/atomix"
	_ "github.com/onosproject/onos-test/test/config"
	_ "github.com/onosproject/onos-test/test/integration"
	_ "github.com/onosproject/onos-test/test/topo"
	test "github.com/onosproject/onos-test/test/kubetest"
)

func main() {
	kubetest.RegisterTests("suite-one", &test.TestsOne{})
	kubetest.RegisterTests("suite-two", &test.TestsTwo{})
	kubetest.RegisterTests("suite-three", &test.TestsThree{})
	kubetest.RegisterBenchmarks("benchmarks-one", &test.BenchmarksOne{})
	kubetest.Main()
}
