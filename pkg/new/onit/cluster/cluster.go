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

package cluster

import (
	atomixcontroller "github.com/atomix/atomix-k8s-controller/pkg/client/clientset/versioned"
	"github.com/onosproject/onos-test/pkg/new/kube"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

// New returns a new onit Env
func New(kube kube.API) *Cluster {
	return &Cluster{
		client: &client{
			namespace:        kube.Namespace(),
			kubeClient:       kubernetes.NewForConfigOrDie(kube.Config()),
			atomixClient:     atomixcontroller.NewForConfigOrDie(kube.Config()),
			extensionsClient: apiextension.NewForConfigOrDie(kube.Config()),
		},
	}
}

// Cluster facilitates modifying subsystems in Kubernetes
type Cluster struct {
	*client
}

// Atomix returns the Atomix service
func (c *Cluster) Atomix() *Atomix {
	return newAtomix(c.client)
}

// Database returns the database service
func (c *Cluster) Database() *Database {
	return newDatabase(c.client)
}

// Topo returns the topo service
func (c *Cluster) Topo() *Topo {
	return newTopo(c.client)
}

// Config returns the configuration service
func (c *Cluster) Config() *Config {
	return newConfig(c.client)
}

// Apps returns the cluster applications
func (c *Cluster) Apps() *Apps {
	return newApps(c.client)
}

// Simulators returns the cluster simulators
func (c *Cluster) Simulators() *Simulators {
	return newSimulators(c.client)
}

// Networks returns the cluster networks
func (c *Cluster) Networks() *Networks {
	return newNetworks(c.client)
}
