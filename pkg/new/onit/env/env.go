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
	"github.com/onosproject/onos-test/pkg/new/kubetest"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

// New returns a new onit Env
func New(kube kubetest.KubeAPI) Env {
	env := &testEnv{
		namespace:        kube.Namespace(),
		kubeClient:       kubernetes.NewForConfigOrDie(kube.Config()),
		atomixClient:     atomixcontroller.NewForConfigOrDie(kube.Config()),
		extensionsclient: apiextension.NewForConfigOrDie(kube.Config()),
	}
	env.atomix = &atomixEnv{
		serviceEnv: &serviceEnv{
			testEnv: env,
		},
	}
	env.database = &databaseEnv{
		testEnv: env,
	}
	env.topo = &topoEnv{
		serviceEnv: &serviceEnv{
			testEnv: env,
		},
	}
	env.config = &configEnv{
		serviceEnv: &serviceEnv{
			testEnv: env,
		},
	}
	env.simulators = &simulatorsEnv{
		testEnv: env,
	}
	env.networks = &networksEnv{
		testEnv: env,
	}
	return env
}

// Env is an interface for tests to operate on the ONOS environment
type Env interface {
	// Atomix returns the Atomix environment
	Atomix() AtomixEnv

	// Database returns the database environment
	Database() DatabaseEnv

	// Topo returns the topo environment
	Topo() TopoEnv

	// Config returns the config environment
	Config() ConfigEnv

	// Simulators returns the simulators environment
	Simulators() SimulatorsEnv

	// Networks returns the networks environment
	Networks() NetworksEnv

	// Apps returns the applications environment
	Apps() AppsEnv
}

// testEnv is an implementation of the Env interface
type testEnv struct {
	namespace        string
	kubeClient       *kubernetes.Clientset
	atomixClient     *atomixcontroller.Clientset
	extensionsclient *apiextension.Clientset
	atomix           *atomixEnv
	database         *databaseEnv
	topo             *topoEnv
	config           *configEnv
	simulators       *simulatorsEnv
	networks         *networksEnv
	apps             *appsEnv
}

func (e *testEnv) Atomix() AtomixEnv {
	return e.atomix
}

func (e *testEnv) Database() DatabaseEnv {
	return e.database
}

func (e *testEnv) Topo() TopoEnv {
	return e.topo
}

func (e *testEnv) Config() ConfigEnv {
	return e.config
}

func (e *testEnv) Simulators() SimulatorsEnv {
	return e.simulators
}

func (e *testEnv) Networks() NetworksEnv {
	return e.networks
}

func (e *testEnv) Apps() AppsEnv {
	return e.apps
}
