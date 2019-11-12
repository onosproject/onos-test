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

package env

import (
	"github.com/onosproject/onos-test/pkg/onit/cluster"
)

// NetworksEnv provides the networks environment
type NetworksEnv interface {
	// List returns a list of networks in the environment
	List() []NetworkEnv

	// Get returns the environment for a network service by name
	Get(name string) NetworkEnv

	// New adds a new network to the environment
	New() NetworkSetup
}

var _ NetworksEnv = &clusterNetworksEnv{}

// clusterNetworksEnv is an implementation of the Networks interface
type clusterNetworksEnv struct {
	networks *cluster.Networks
}

func (e *clusterNetworksEnv) List() []NetworkEnv {
	clusterNetworks := e.networks.List()
	networks := make([]NetworkEnv, len(clusterNetworks))
	for i, network := range clusterNetworks {
		networks[i] = e.Get(network.Name())
	}
	return networks
}

func (e *clusterNetworksEnv) Get(name string) NetworkEnv {
	network := e.networks.Get(name)
	return &clusterNetworkEnv{
		clusterNodeEnv: &clusterNodeEnv{
			node: network.Node,
		},
		network: network,
	}
}

func (e *clusterNetworksEnv) New() NetworkSetup {
	return &clusterNetworkSetup{
		network: e.networks.New(),
	}
}
