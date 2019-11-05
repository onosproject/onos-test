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
	"github.com/onosproject/onos-test/pkg/new/kube"
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
)

// New returns a new onit Deployment
func New(kube kube.API) Deployment {
	return &clusterDeployment{
		cluster: cluster.New(kube),
	}
}

// Deploy is an interface for deploying ONOS subsystems at runtime
type Deploy interface {
	// Deploy deploys the service
	Deploy() error

	// DeployOrDie deploys the service and panics if the deployment fails
	DeployOrDie()
}

// Deployment is an interface for deploying subsystems at runtime
type Deployment interface {
	// App returns a new app deployer
	App(name string) App

	// Simulator returns a new simulator deployer
	Simulator(name string) Simulator

	// Network returns a new network deployer
	Network(name string) Network
}

// clusterDeployment is an implementation of the Deployment interface
type clusterDeployment struct {
	cluster *cluster.Cluster
}

func (d *clusterDeployment) App(name string) App {
	app := d.cluster.Apps().Get(name)
	deploy := &clusterApp{
		clusterServiceType: &clusterServiceType{
			clusterService: &clusterService{
				service: app.Service,
			},
		},
		app: app,
	}
	deploy.clusterServiceType.deploy = deploy
	return deploy
}

func (d *clusterDeployment) Simulator(name string) Simulator {
	simulator := d.cluster.Simulators().Get(name)
	deploy := &clusterSimulator{
		clusterNodeType: &clusterNodeType{
			clusterNode: &clusterNode{
				node: simulator.Node,
			},
		},
		simulator: simulator,
	}
	deploy.clusterNodeType.deploy = deploy
	return deploy
}

func (d *clusterDeployment) Network(name string) Network {
	network := d.cluster.Networks().Get(name)
	deploy := &clusterNetwork{
		clusterNodeType: &clusterNodeType{
			clusterNode: &clusterNode{
				node: network.Node,
			},
		},
		network: network,
	}
	deploy.clusterNodeType.deploy = deploy
	return deploy
}
