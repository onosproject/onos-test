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

package cli

const (
	atomixService = "atomix"
	raftService   = "raft"
	configService = "config"
	topoService   = "topo"
)

const (
	defaultAtomixImage    = "atomix/atomix-k8s-controller:latest"
	defaultRaftImage      = "atomix/atomix-raft-node:latest"
	defaultConfigImage    = "onosproject/onos-config:latest"
	defaultTopoImage      = "onosproject/onos-topo:latest"
	defaultMininetImage   = "opennetworkinglab/mininet:latest"
	defaultSimulatorImage = "onosproject/simulators:latest"
)
