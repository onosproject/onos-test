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
	"github.com/onosproject/onos-test/pkg/new/kubetest"
	"github.com/onosproject/onos-test/pkg/new/util/logging"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"path/filepath"
	"runtime"
)

var (
	_, path, _, _ = runtime.Caller(0)
	certsPath     = filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(path)))), "certs")
)

// New returns a new onit Env
func New(kube kube.API) *Cluster {
	client := &client{
		namespace:        kube.Namespace(),
		kubeClient:       kubernetes.NewForConfigOrDie(kube.Config()),
		atomixClient:     atomixcontroller.NewForConfigOrDie(kube.Config()),
		extensionsClient: apiextension.NewForConfigOrDie(kube.Config()),
	}
	return &Cluster{
		client:     client,
		atomix:     newAtomix(client),
		database:   newDatabase(client),
		topo:       newTopo(client),
		config:     newConfig(client),
		apps:       newApps(client),
		simulators: newSimulators(client),
		networks:   newNetworks(client),
	}
}

// Cluster facilitates modifying subsystems in Kubernetes
type Cluster struct {
	*client
	atomix     *Atomix
	database   *Database
	topo       *Topo
	config     *Config
	apps       *Apps
	simulators *Simulators
	networks   *Networks
}

// Atomix returns the Atomix service
func (c *Cluster) Atomix() *Atomix {
	return c.atomix
}

// Database returns the database service
func (c *Cluster) Database() *Database {
	return c.database
}

// Topo returns the topo service
func (c *Cluster) Topo() *Topo {
	return c.topo
}

// Config returns the configuration service
func (c *Cluster) Config() *Config {
	return c.config
}

// Apps returns the cluster applications
func (c *Cluster) Apps() *Apps {
	return c.apps
}

// Simulators returns the cluster simulators
func (c *Cluster) Simulators() *Simulators {
	return c.simulators
}

// Networks returns the cluster networks
func (c *Cluster) Networks() *Networks {
	return c.networks
}

// Create creates the cluster
func (c *Cluster) Create() error {
	step := logging.NewStep(c.namespace, "Setup ONOS cluster")
	step.Start()
	cluster := kubetest.NewTestCluster(c.namespace)
	if err := cluster.Create(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// Delete deletes the cluster
func (c *Cluster) Delete() error {
	step := logging.NewStep(c.namespace, "Tear down ONOS cluster")
	step.Start()
	cluster := kubetest.NewTestCluster(c.namespace)
	if err := cluster.Delete(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}
