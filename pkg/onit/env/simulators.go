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

// SimulatorsEnv provides the simulators environment
type SimulatorsEnv interface {
	// List returns a list of simulators in the environment
	List() []SimulatorEnv

	// Get returns the environment for a simulator service by name
	Get(name string) SimulatorEnv

	// New adds a new simulator to the environment
	New() SimulatorSetup
}

var _ SimulatorsEnv = &clusterSimulatorsEnv{}

// clusterSimulatorsEnv is an implementation of the Simulators interface
type clusterSimulatorsEnv struct {
	simulators *cluster.Simulators
}

func (e *clusterSimulatorsEnv) List() []SimulatorEnv {
	clusterNetworks := e.simulators.List()
	networks := make([]SimulatorEnv, len(clusterNetworks))
	for i, network := range clusterNetworks {
		networks[i] = e.Get(network.Name())
	}
	return networks
}

func (e *clusterSimulatorsEnv) Get(name string) SimulatorEnv {
	simulator := e.simulators.Get(name)
	return &clusterSimulatorEnv{
		clusterNodeEnv: &clusterNodeEnv{
			node: simulator.Node,
		},
		simulator: simulator,
	}
}

func (e *clusterSimulatorsEnv) New() SimulatorSetup {
	return &clusterSimulatorSetup{
		simulator: e.simulators.New(),
	}
}
