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
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
)

// Config is an interface for setting up config nodes
type Config interface {
	ServiceType
	concurrentSetup

	// Nodes sets the number of nodes to deploy
	Nodes(nodes int) Config
}

var _ Config = &clusterConfig{}

// clusterConfig is an implementation of the Config interface
type clusterConfig struct {
	*clusterServiceType
	config *cluster.Config
}

func (s *clusterConfig) Nodes(nodes int) Config {
	s.config.SetReplicas(nodes)
	return s
}

func (s *clusterConfig) create() error {
	return s.config.Create()
}

func (s *clusterConfig) waitForStart() error {
	return s.config.AwaitReady()
}
