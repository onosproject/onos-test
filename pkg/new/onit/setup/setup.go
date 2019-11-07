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
	"github.com/onosproject/onos-test/pkg/new/kube"
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
)

const raftGroup = "raft"

// New returns a new onit Setup
func New(kube kube.API) Setup {
	return &clusterSetup{
		cluster: cluster.New(kube),
	}
}

// Setup is an interface for setting up ONOS clusters
type Setup interface {
	// Atomix returns the setup configuration for the Atomix controller
	Atomix() Atomix

	// Database returns the setup configuration for the key-value store
	Database() Database

	// Topo returns the setup configuration for the ONOS topo service
	Topo() Topo

	// Config returns the setup configuration for the ONOS config service
	Config() Config

	// Setup sets up the cluster
	Setup() error

	// SetupOrDie sets up the cluster and panics if the setup fails
	SetupOrDie()
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

// clusterSetup is an implementation of the Setup interface
type clusterSetup struct {
	cluster *cluster.Cluster
}

func (s *clusterSetup) Atomix() Atomix {
	return &clusterAtomix{
		atomix: s.cluster.Atomix(),
	}
}

func (s *clusterSetup) Database() Database {
	return &clusterDatabase{
		group: s.cluster.Database().Partitions(raftGroup),
	}
}

func (s *clusterSetup) Topo() Topo {
	return &clusterTopo{
		topo: s.cluster.Topo(),
	}
}

func (s *clusterSetup) Config() Config {
	return &clusterConfig{
		config: s.cluster.Config(),
	}
}

func (s *clusterSetup) Setup() error {
	// Set up the Atomix controller
	if err := s.Atomix().(sequentialSetup).setup(); err != nil {
		return err
	}

	// Create the database and services concurrently
	if err := s.Database().(concurrentSetup).create(); err != nil {
		return err
	}
	if err := s.Topo().(concurrentSetup).create(); err != nil {
		return err
	}
	if err := s.Config().(concurrentSetup).create(); err != nil {
		return err
	}

	// Wait for the database and services to start up
	if err := s.Database().(concurrentSetup).waitForStart(); err != nil {
		return err
	}
	if err := s.Topo().(concurrentSetup).waitForStart(); err != nil {
		return err
	}
	if err := s.Config().(concurrentSetup).waitForStart(); err != nil {
		return err
	}
	return nil
}

func (s *clusterSetup) SetupOrDie() {
	if err := s.Setup(); err != nil {
		panic(err)
	}
}

var _ Setup = &clusterSetup{}
