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

// Simulators provides the simulators environment
type Simulators interface {
	// List returns a list of simulators in the environment
	List() []Simulator

	// Get returns the environment for a simulator service by name
	Get(name string) Simulator

	// Add adds a new simulator to the environment
	Add(name string) deploy.Simulator
}

var _ Simulators = &clusterSimulators{}

// clusterSimulators is an implementation of the Simulators interface
type clusterSimulators struct {
	deployment deploy.Deployment
	simulators *cluster.Simulators
}

func (e *clusterSimulators) List() []Simulator {
	clusterNetworks := e.simulators.List()
	networks := make([]Simulator, len(clusterNetworks))
	for i, network := range clusterNetworks {
		networks[i] = e.Get(network.Name())
	}
	return networks
}

func (e *clusterSimulators) Get(name string) Simulator {
	simulator := e.simulators.Get(name)
	return &clusterSimulator{
		clusterNode: &clusterNode{
			node: simulator.Node,
		},
		simulator: simulator,
	}
}

func (e *clusterSimulators) Add(name string) deploy.Simulator {
	return e.deployment.Simulator(name)
}
