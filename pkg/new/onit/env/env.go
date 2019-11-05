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
	atomixcontroller "github.com/atomix/atomix-k8s-controller/pkg/client/clientset/versioned"
	"github.com/onosproject/onos-test/pkg/new/kube"
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
	"github.com/onosproject/onos-test/pkg/new/onit/deploy"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

// New returns a new onit Env
func New(kube kube.API) Env {
	cluster := cluster.New(kube)
	atomix := cluster.Atomix()
	group := cluster.Database().Partitions("raft")
	topo := cluster.Topo()
	config := cluster.Config()
	apps := cluster.Apps()
	simulators := cluster.Simulators()
	networks := cluster.Networks()
	deployment := deploy.New(kube)
	return &clusterEnv{
		namespace:        kube.Namespace(),
		kubeClient:       kubernetes.NewForConfigOrDie(kube.Config()),
		atomixClient:     atomixcontroller.NewForConfigOrDie(kube.Config()),
		extensionsclient: apiextension.NewForConfigOrDie(kube.Config()),
		deployment:       deploy.New(kube),
		atomix: &clusterAtomix{
			clusterService: &clusterService{
				service: atomix.Service,
			},
		},
		database: &clusterDatabase{
			group: group,
		},
		topo: &clusterTopo{
			clusterService: &clusterService{
				service: topo.Service,
			},
		},
		config: &clusterConfig{
			clusterService: &clusterService{
				service: config.Service,
			},
		},
		apps: &clusterApps{
			deployment: deployment,
			apps:       apps,
		},
		simulators: &clusterSimulators{
			deployment: deployment,
			simulators: simulators,
		},
		networks: &clusterNetworks{
			deployment: deployment,
			networks:   networks,
		},
	}
}

// Env is an interface for tests to operate on the ONOS environment
type Env interface {
	// Atomix returns the Atomix environment
	Atomix() Atomix

	// Database returns the database environment
	Database() Database

	// Topo returns the topo environment
	Topo() Topo

	// Config returns the config environment
	Config() Config

	// Simulators returns the simulators environment
	Simulators() Simulators

	// Simulator returns the environment for a simulator by name
	Simulator(name string) Simulator

	// AddSimulator returns a Simulator deployment for adding a simulator to the cluster
	AddSimulator(name string) deploy.Simulator

	// Networks returns the networks environment
	Networks() Networks

	// Network returns the environment for a network by name
	Network(name string) Network

	// AddNetwork returns a Network deployment for adding a network to the cluster
	AddNetwork(name string) deploy.Network

	// Apps returns the applications environment
	Apps() Apps

	// App returns the environment for an app by name
	App(name string) App

	// AddApp returns an App deployment for adding an application to the cluster
	AddApp(name string) deploy.App
}

// clusterEnv is an implementation of the Env interface
type clusterEnv struct {
	namespace        string
	kubeClient       *kubernetes.Clientset
	atomixClient     *atomixcontroller.Clientset
	extensionsclient *apiextension.Clientset
	deployment       deploy.Deployment
	atomix           *clusterAtomix
	database         *clusterDatabase
	topo             *clusterTopo
	config           *clusterConfig
	simulators       *clusterSimulators
	networks         *clusterNetworks
	apps             *clusterApps
}

func (e *clusterEnv) Atomix() Atomix {
	return e.atomix
}

func (e *clusterEnv) Database() Database {
	return e.database
}

func (e *clusterEnv) Topo() Topo {
	return e.topo
}

func (e *clusterEnv) Config() Config {
	return e.config
}

func (e *clusterEnv) Simulators() Simulators {
	return e.simulators
}

func (e *clusterEnv) Simulator(name string) Simulator {
	return e.Simulators().Get(name)
}

func (e *clusterEnv) AddSimulator(name string) deploy.Simulator {
	return e.Simulators().Add(name)
}

func (e *clusterEnv) Networks() Networks {
	return e.networks
}

func (e *clusterEnv) Network(name string) Network {
	return e.Networks().Get(name)
}

func (e *clusterEnv) AddNetwork(name string) deploy.Network {
	return e.Networks().Add(name)
}

func (e *clusterEnv) Apps() Apps {
	return e.apps
}

func (e *clusterEnv) App(name string) App {
	return e.Apps().Get(name)
}

func (e *clusterEnv) AddApp(name string) deploy.App {
	return e.Apps().Add(name)
}
