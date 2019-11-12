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

package cluster

type serviceType string

func (s serviceType) name() string {
	return string(s)
}

const (
	atomixType    serviceType = "atomix"
	databaseType  serviceType = "database"
	partitionType serviceType = "partition"
	cliType       serviceType = "cli"
	topoType      serviceType = "topo"
	configType    serviceType = "config"
	appType       serviceType = "app"
	simulatorType serviceType = "simulator"
	networkType   serviceType = "network"
)

const (
	atomixImage    = "atomix/atomix-k8s-controller:latest"
	raftImage      = "atomix/atomix-raft-node:latest"
	cliImage       = "onosproject/onos-cli:latest"
	topoImage      = "onosproject/onos-topo:latest"
	configImage    = "onosproject/onos-config:latest"
	simulatorImage = "onosproject/device-simulator:latest"
	networkImage   = "opennetworking/mn-stratum:latest"
)

const (
	typeLabel = "type"
)
