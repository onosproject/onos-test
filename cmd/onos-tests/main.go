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
	"github.com/onosproject/onos-test/pkg/onit"
	"github.com/onosproject/onos-test/test/atomix"
	"github.com/onosproject/onos-test/test/config"
	"github.com/onosproject/onos-test/test/topo"
)

func main() {
	onit.RegisterTests("atomix", &atomix.SmokeTestSuite{})
	onit.RegisterTests("atomix-ha", &atomix.HATestSuite{})
	onit.RegisterTests("topo", &topo.TestSuite{})
	onit.RegisterTests("config", &config.SmokeTestSuite{})
	onit.RegisterTests("config-cli", &config.CLITestSuite{})
	onit.RegisterTests("config-ha", &config.HATestSuite{})

	onit.RegisterBenchmarks("atomix", &atomix.BenchmarkSuite{})
	onit.RegisterBenchmarks("topo", &topo.BenchmarkSuite{})

	onit.Main()
}
