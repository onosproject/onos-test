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
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/onit/cluster"
)

// New returns a new onit ClusterEnv
func New(kube kube.API) ClusterEnv {
	return &clusterEnv{
		cluster: cluster.New(kube),
	}
}

var env ClusterEnv

// getEnv gets the current environment
func getEnv() ClusterEnv {
	if env == nil {
		env = New(kube.GetAPIFromEnvOrDie())
	}
	return env
}

// Atomix returns the Atomix environment
func Atomix() AtomixEnv {
	return getEnv().Atomix()
}

// Storage returns the database environment
func Storage() StorageEnv {
	return getEnv().Storage()
}

// CLI returns the CLI environment
func CLI() CLIEnv {
	return getEnv().CLI()
}

// Topo returns the onos-topo environment
func Topo() TopoEnv {
	return getEnv().Topo()
}

// RAN returns the onos-ric environment
func RAN() RICEnv {
	return getEnv().RIC()
}

// Config returns the onos-config environment
func Config() ConfigEnv {
	return getEnv().Config()
}

// Simulators returns the device simulators environment
func Simulators() SimulatorsEnv {
	return getEnv().Simulators()
}

// Simulator returns the environment for a device simulator
func Simulator(name string) SimulatorEnv {
	return getEnv().Simulator(name)
}

// NewSimulator returns the setup configuration for a new device simulator
func NewSimulator() SimulatorSetup {
	return getEnv().NewSimulator()
}

// AddSimulators returns a new SimulatorsSetup for adding multiple simulators concurrently
func AddSimulators(simulators ...SimulatorSetup) SimulatorsSetup {
	return getEnv().AddSimulators(simulators...)
}

// Networks returns the networks environment
func Networks() NetworksEnv {
	return getEnv().Networks()
}

// Network returns the environment for a network
func Network(name string) NetworkEnv {
	return getEnv().Network(name)
}

// History returns the environment for the history
func History() HistoryEnv {
	return getEnv().History()
}

// NewNetwork returns the setup configuration for a new network
func NewNetwork() NetworkSetup {
	return getEnv().NewNetwork()
}

// Apps returns the environment for applications
func Apps() AppsEnv {
	return getEnv().Apps()
}

// App returns the environment for an application
func App(name string) AppEnv {
	return getEnv().App(name)
}

// NewApp returns the setup configuration for a new application
func NewApp() AppSetup {
	return getEnv().NewApp()
}

// ClusterEnv is an interface for tests to operate on the ONOS environment
type ClusterEnv interface {
	// Atomix returns the Atomix environment
	Atomix() AtomixEnv

	// Storage returns the storage environment
	Storage() StorageEnv

	// CLI returns the CLI environment
	CLI() CLIEnv

	// Topo returns the topo environment
	Topo() TopoEnv

	// RIC returns the RAN environment
	RIC() RICEnv

	// Config returns the config environment
	Config() ConfigEnv

	// Simulators returns the simulators environment
	Simulators() SimulatorsEnv

	// Simulator returns the environment for a simulator by name
	Simulator(name string) SimulatorEnv

	// NewSimulator returns a new SimulatorSetup for adding a simulator to the cluster
	NewSimulator() SimulatorSetup

	// AddSimulators returns a new SimulatorsSetup for adding multiple simulators concurrently
	AddSimulators(simulators ...SimulatorSetup) SimulatorsSetup

	// Networks returns the networks environment
	Networks() NetworksEnv

	// Network returns the environment for a network by name
	Network(name string) NetworkEnv

	// History returns the environment for the history
	History() HistoryEnv

	// NewNetwork returns a new NetworkSetup for adding a network to the cluster
	NewNetwork() NetworkSetup

	// Apps returns the applications environment
	Apps() AppsEnv

	// App returns the environment for an app by name
	App(name string) AppEnv

	// NewApp returns a new AppSetup for adding an application to the cluster
	NewApp() AppSetup
}

// clusterEnv is an implementation of the Env interface
type clusterEnv struct {
	cluster *cluster.Cluster
}

func (e *clusterEnv) Atomix() AtomixEnv {
	return &clusterAtomixEnv{
		clusterServiceEnv: &clusterServiceEnv{
			clusterDeploymentEnv: &clusterDeploymentEnv{
				deployment: e.cluster.Atomix().Deployment,
			},
		},
	}
}

func (e *clusterEnv) Storage() StorageEnv {
	return &clusterStorageEnv{
		database: e.cluster.Database(),
	}
}

func (e *clusterEnv) CLI() CLIEnv {
	return &clusterCLIEnv{
		clusterDeploymentEnv: &clusterDeploymentEnv{
			deployment: e.cluster.CLI().Deployment,
		},
	}
}

func (e *clusterEnv) Topo() TopoEnv {
	return &clusterTopoEnv{
		clusterServiceEnv: &clusterServiceEnv{
			clusterDeploymentEnv: &clusterDeploymentEnv{
				deployment: e.cluster.Topo().Deployment,
			},
			service: e.cluster.Topo().Service,
		},
	}
}

func (e *clusterEnv) RIC() RICEnv {
	return &clusterRICEnv{
		clusterServiceEnv: &clusterServiceEnv{
			clusterDeploymentEnv: &clusterDeploymentEnv{
				deployment: e.cluster.RAN().Deployment,
			},
			service: e.cluster.RAN().Service,
		},
	}
}

func (e *clusterEnv) Config() ConfigEnv {
	return &clusterConfigEnv{
		clusterServiceEnv: &clusterServiceEnv{
			clusterDeploymentEnv: &clusterDeploymentEnv{
				deployment: e.cluster.Topo().Deployment,
			},
			service: e.cluster.Config().Service,
		},
	}
}

func (e *clusterEnv) Simulators() SimulatorsEnv {
	return &clusterSimulatorsEnv{
		simulators: e.cluster.Simulators(),
	}
}

func (e *clusterEnv) Simulator(name string) SimulatorEnv {
	return e.Simulators().Get(name)
}

func (e *clusterEnv) NewSimulator() SimulatorSetup {
	return e.Simulators().New()
}

func (e *clusterEnv) AddSimulators(simulators ...SimulatorSetup) SimulatorsSetup {
	return &clusterSimulatorsSetup{
		simulators: e.cluster.Simulators(),
		setups:     make([]SimulatorSetup, 0),
	}
}

func (e *clusterEnv) Networks() NetworksEnv {
	return &clusterNetworksEnv{
		networks: e.cluster.Networks(),
	}
}

func (e *clusterEnv) History() HistoryEnv {
	return &clusterHistoryEnv{
		history: e.cluster.History(),
	}
}

func (e *clusterEnv) Network(name string) NetworkEnv {
	return e.Networks().Get(name)
}

func (e *clusterEnv) NewNetwork() NetworkSetup {
	return e.Networks().New()
}

func (e *clusterEnv) Apps() AppsEnv {
	return &clusterAppsEnv{
		apps: e.cluster.Apps(),
	}
}

func (e *clusterEnv) App(name string) AppEnv {
	return e.Apps().Get(name)
}

func (e *clusterEnv) NewApp() AppSetup {
	return e.Apps().New()
}
