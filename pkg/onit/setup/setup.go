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
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"sync"
)

const raftGroup = "raft"

// New returns a new onit ClusterSetup
func New(kube kube.API) ClusterSetup {
	return &clusterSetup{
		cluster: cluster.New(kube),
		apps:    make(map[string]AppSetup),
	}
}

var setup ClusterSetup

// getSetup gets the current setup
func getSetup() ClusterSetup {
	if setup == nil {
		setup = New(kube.GetAPIFromEnvOrDie())
	}
	return setup
}

// Atomix returns the setup configuration for the Atomix controller
func Atomix() AtomixSetup {
	return getSetup().Atomix()
}

// Database returns the setup configuration for the key-value store
func Database() DatabaseSetup {
	return getSetup().Database()
}

// CLI returns the setup configuration for the CLI service
func CLI() CLISetup {
	return getSetup().CLI()
}

// App returns the setup configuration for an application
func App(name string) AppSetup {
	return getSetup().App(name)
}

// Topo returns the setup configuration for the topo service
func Topo() TopoSetup {
	return getSetup().Topo()
}

// Config returns the setup configuration for the config service
func Config() ConfigSetup {
	return getSetup().Config()
}

// Setup sets up the cluster
func Setup() error {
	return getSetup().Setup()
}

// SetupOrDie sets up the cluster and panics if setup fails
func SetupOrDie() { //nolint:golint
	getSetup().SetupOrDie()
}

// ClusterSetup is an interface for setting up ONOS clusters
type ClusterSetup interface {
	// Atomix returns the setup configuration for the Atomix controller
	Atomix() AtomixSetup

	// Database returns the setup configuration for the key-value store
	Database() DatabaseSetup

	// CLI returns the setup configuration for the ONSO CLI service
	CLI() CLISetup

	// Topo returns the setup configuration for the ONOS topo service
	Topo() TopoSetup

	// Config returns the setup configuration for the ONOS config service
	Config() ConfigSetup

	// App returns the setup configuration for an application
	App(name string) AppSetup

	// Setup sets up the cluster
	Setup() error

	// SetupOrDie sets up the cluster and panics if the setup fails
	SetupOrDie()
}

// serviceSetup is a setup step for a single service
type serviceSetup interface {
	setup() error
}

// clusterSetup is an implementation of the Setup interface
type clusterSetup struct {
	cluster *cluster.Cluster
	apps    map[string]AppSetup
}

func (s *clusterSetup) Atomix() AtomixSetup {
	return &clusterAtomixSetup{
		atomix: s.cluster.Atomix(),
	}
}

func (s *clusterSetup) Database() DatabaseSetup {
	return &clusterDatabaseSetup{
		group: s.cluster.Database().Partitions(raftGroup),
	}
}

func (s *clusterSetup) CLI() CLISetup {
	return &clusterCLISetup{
		cli: s.cluster.CLI(),
	}
}

func (s *clusterSetup) Topo() TopoSetup {
	return &clusterTopoSetup{
		topo: s.cluster.Topo(),
	}
}

func (s *clusterSetup) Config() ConfigSetup {
	return &clusterConfigSetup{
		config: s.cluster.Config(),
	}
}

func (s *clusterSetup) App(name string) AppSetup {
	if app, ok := s.apps[name]; ok {
		return app
	}

	app := s.cluster.Apps().New()
	app.SetName(name)
	setup := &clusterAppSetup{
		app: app,
	}
	s.apps[name] = setup
	return setup
}

func (s *clusterSetup) Setup() error {
	// Set up the Atomix controller
	if err := s.Atomix().(serviceSetup).setup(); err != nil {
		return err
	}

	// Create the database and services concurrently
	wg := &sync.WaitGroup{}
	errCh := make(chan error)

	setupService(s.Database().(serviceSetup), wg, errCh)
	setupService(s.CLI().(serviceSetup), wg, errCh)
	setupService(s.Topo().(serviceSetup), wg, errCh)
	setupService(s.Config().(serviceSetup), wg, errCh)
	for _, app := range s.apps {
		setupService(app.(serviceSetup), wg, errCh)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		return err
	}
	return nil
}

func setupService(setup serviceSetup, wg *sync.WaitGroup, errCh chan<- error) {
	wg.Add(1)
	go func() {
		if err := setup.setup(); err != nil {
			errCh <- err
		}
		wg.Done()
	}()
}

func (s *clusterSetup) SetupOrDie() {
	if err := s.Setup(); err != nil {
		panic(err)
	}
}

var _ ClusterSetup = &clusterSetup{}
