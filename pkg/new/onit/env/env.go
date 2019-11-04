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
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

// New returns a new onit Env
func New(kube kube.API) Env {
	env := &testEnv{
		namespace:        kube.Namespace(),
		kubeClient:       kubernetes.NewForConfigOrDie(kube.Config()),
		atomixClient:     atomixcontroller.NewForConfigOrDie(kube.Config()),
		extensionsclient: apiextension.NewForConfigOrDie(kube.Config()),
	}
	env.atomix = &atomixEnv{
		service: &service{
			testEnv: env,
		},
	}
	env.database = &database{
		testEnv: env,
	}
	env.topo = &topo{
		service: &service{
			testEnv: env,
		},
	}
	env.config = &config{
		service: &service{
			testEnv: env,
		},
	}
	env.simulators = &simulators{
		testEnv: env,
	}
	env.networks = &networks{
		testEnv: env,
	}
	return env
}

// Env is an interface for tests to operate on the ONOS environment
type Env interface {
	// Atomix returns the Atomix environment
	Atomix() AtomixEnv

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

	// Networks returns the networks environment
	Networks() Networks

	// Network returns the environment for a network by name
	Network(name string) Network

	// Apps returns the applications environment
	Apps() Apps

	// App returns the environment for an app by name
	App(name string) App
}

// testEnv is an implementation of the Env interface
type testEnv struct {
	namespace        string
	kubeClient       *kubernetes.Clientset
	atomixClient     *atomixcontroller.Clientset
	extensionsclient *apiextension.Clientset
	atomix           *atomixEnv
	database         *database
	topo             *topo
	config           *config
	simulators       *simulators
	networks         *networks
	apps             *apps
}

func (e *testEnv) Atomix() AtomixEnv {
	return e.atomix
}

func (e *testEnv) Database() Database {
	return e.database
}

func (e *testEnv) Topo() Topo {
	return e.topo
}

func (e *testEnv) Config() Config {
	return e.config
}

func (e *testEnv) Simulators() Simulators {
	return e.simulators
}

func (e *testEnv) Simulator(name string) Simulator {
	return e.Simulators().Get(name)
}

func (e *testEnv) Networks() Networks {
	return e.networks
}

func (e *testEnv) Network(name string) Network {
	return e.Networks().Get(name)
}

func (e *testEnv) Apps() Apps {
	return e.apps
}

func (e *testEnv) App(name string) App {
	return e.Apps().Get(name)
}
