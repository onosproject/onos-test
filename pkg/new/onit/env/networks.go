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
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
	"github.com/onosproject/onos-test/pkg/new/onit/deploy"
)

// Networks provides the networks environment
type Networks interface {
	// List returns a list of networks in the environment
	List() []Network

	// Get returns the environment for a network service by name
	Get(name string) Network

	// Add adds a new network to the environment
	Add(name string) deploy.Network
}

var _ Networks = &clusterNetworks{}

// clusterNetworks is an implementation of the Networks interface
type clusterNetworks struct {
	deployment deploy.Deployment
	networks   *cluster.Networks
}

func (e *clusterNetworks) List() []Network {
	clusterNetworks := e.networks.List()
	networks := make([]Network, len(clusterNetworks))
	for i, network := range clusterNetworks {
		networks[i] = e.Get(network.Name())
	}
	return networks
}

func (e *clusterNetworks) Get(name string) Network {
	network := e.networks.Get(name)
	return &clusterNetwork{
		clusterNode: &clusterNode{
			node: network.Node,
		},
		network: network,
	}
}

func (e *clusterNetworks) Add(name string) deploy.Network {
	return e.deployment.Network(name)
}
