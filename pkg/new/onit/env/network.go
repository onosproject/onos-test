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
)

// Network provides the environment for a network node
type Network interface {
	Node

	// Devices returns a list of devices in the network
	Devices() []Node

	// Remove removes the network
	Remove() error

	// RemoveOrDie removes the network and panics if the remove fails
	RemoveOrDie()
}

var _ Network = &clusterNetwork{}

// clusterService is an implementation of the Network interface
type clusterNetwork struct {
	*clusterNode
	network *cluster.Network
}

func (e *clusterNetwork) Devices() []Node {
	clusterDevices := e.network.Devices()
	devices := make([]Node, len(clusterDevices))
	for i, node := range clusterDevices {
		devices[i] = &clusterNode{
			node: node,
		}
	}
	return devices
}

func (e *clusterNetwork) Remove() error {
	return e.network.Remove()
}

func (e *clusterNetwork) RemoveOrDie() {
	if err := e.Remove(); err != nil {
		panic(err)
	}
}
