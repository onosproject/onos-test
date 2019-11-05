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

package deploy

import (
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
)

// Network is an interface for deploying up a network
type Network interface {
	Deploy
	NodeType

	// Single creates a single node topology
	Single() Network

	// Linear creates a linear topology with the given number of devices
	Linear(devices int) Network
}

var _ Network = &clusterNetwork{}

// clusterNetwork is an implementation of the Network interface
type clusterNetwork struct {
	*clusterNodeType
	network *cluster.Network
}

func (s *clusterNetwork) Single() Network {
	s.network.SetSingle()
	return s
}

func (s *clusterNetwork) Linear(devices int) Network {
	s.network.SetLinear(devices)
	return s
}

func (s *clusterNetwork) Deploy() error {
	return s.network.Add()
}

func (s *clusterNetwork) DeployOrDie() {
	if err := s.Deploy(); err != nil {
		panic(err)
	}
}
