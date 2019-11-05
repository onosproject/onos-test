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
	atomixcontroller "github.com/atomix/atomix-k8s-controller/pkg/client/clientset/versioned"
	"github.com/onosproject/onos-test/pkg/new/kube"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

// New returns a new onit Deployment
func New(kube kube.API) Deployment {
	return &testDeployment{
		namespace:        kube.Namespace(),
		kubeClient:       kubernetes.NewForConfigOrDie(kube.Config()),
		atomixClient:     atomixcontroller.NewForConfigOrDie(kube.Config()),
		extensionsclient: apiextension.NewForConfigOrDie(kube.Config()),
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

// testDeployment is an implementation of the Deployment interface
type testDeployment struct {
	namespace        string
	kubeClient       *kubernetes.Clientset
	atomixClient     *atomixcontroller.Clientset
	extensionsclient *apiextension.Clientset
}

func (d *testDeployment) App(name string) App {
	return newAppDeploy(name, d)
}

func (d *testDeployment) Simulator(name string) Simulator {
	return newSimulatorDeploy(name, d)
}

func (d *testDeployment) Network(name string) Network {
	return newNetworkDeploy(name, d)
}
