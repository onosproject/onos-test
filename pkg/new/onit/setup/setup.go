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

package setup

import (
	atomixcontroller "github.com/atomix/atomix-k8s-controller/pkg/client/clientset/versioned"
	"github.com/onosproject/onos-test/pkg/new/kubetest"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

// New returns a new onit Setup
func New(kube kubetest.KubeAPI) TestSetup {
	setup := &testSetup{
		namespace:        kube.Namespace(),
		kubeClient:       kubernetes.NewForConfigOrDie(kube.Config()),
		atomixClient:     atomixcontroller.NewForConfigOrDie(kube.Config()),
		extensionsclient: apiextension.NewForConfigOrDie(kube.Config()),
	}
	setup.atomix = &atomixSetup{
		testSetup: setup,
	}
	setup.database = &databaseSetup{
		testSetup: setup,
	}
	setup.topo = &topoSetup{
		testSetup: setup,
	}
	setup.config = &configSetup{
		testSetup: setup,
	}
	return setup
}

// Setup is an interface for setting up ONOS clusters
type Setup interface {
	Setup() error
}

// TestSetup is an interface for setting up ONOS clusters
type TestSetup interface {
	Setup
	Atomix() AtomixSetup
	Database() DatabaseSetup
	Topo() TopoSetup
	Config() ConfigSetup
}

// sequentialSetup is a setup step that must run sequentially
type sequentialSetup interface {
	setup() error
}

// concurrentSetup is a setup step that can run concurrently with other steps
type concurrentSetup interface {
	create() error
	waitForStart() error
}

// testSetup is an implementation of the Setup interface
type testSetup struct {
	namespace        string
	kubeClient       *kubernetes.Clientset
	atomixClient     *atomixcontroller.Clientset
	extensionsclient *apiextension.Clientset
	atomix           AtomixSetup
	database         DatabaseSetup
	topo             TopoSetup
	config           ConfigSetup
}

func (s *testSetup) Atomix() AtomixSetup {
	return s.atomix
}

func (s *testSetup) Database() DatabaseSetup {
	return s.database
}

func (s *testSetup) Topo() TopoSetup {
	return s.topo
}

func (s *testSetup) Config() ConfigSetup {
	return s.config
}

func (s *testSetup) Setup() error {
	// Set up the Atomix controller
	if err := s.atomix.setup(); err != nil {
		return err
	}

	// Create the database and services concurrently
	if err := s.database.create(); err != nil {
		return err
	}
	if err := s.topo.create(); err != nil {
		return err
	}
	if err := s.config.create(); err != nil {
		return err
	}

	// Wait for the database and services to start up
	if err := s.database.waitForStart(); err != nil {
		return err
	}
	if err := s.topo.waitForStart(); err != nil {
		return err
	}
	if err := s.config.waitForStart(); err != nil {
		return err
	}
	return nil
}

var _ Setup = &testSetup{}
