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
	"github.com/onosproject/onos-test/pkg/test"
	"github.com/onosproject/onos-test/test/config"
	"github.com/onosproject/onos-test/test/grpc"
	"github.com/onosproject/onos-test/test/nopaxos"
	"github.com/onosproject/onos-test/test/raft"
	"github.com/onosproject/onos-test/test/topo"
)

func main() {
	test.RegisterTests("raft", &raft.SmokeTestSuite{})
	test.RegisterTests("raft-ha", &raft.HATestSuite{})
	test.RegisterTests("nopaxos", &nopaxos.SmokeTestSuite{})
	test.RegisterTests("nopaxos-ha", &nopaxos.HATestSuite{})
	test.RegisterTests("topo", &topo.TestSuite{})
	test.RegisterTests("config", &config.SmokeTestSuite{})
	test.RegisterTests("config-cli", &config.CLITestSuite{})
	test.RegisterTests("config-ha", &config.HATestSuite{})

	test.RegisterBenchmarks("nopaxos", &nopaxos.MapBenchmarkSuite{})
	test.RegisterBenchmarks("grpc", &grpc.GRPCBenchmarkSuite{})

	test.Main()
}
