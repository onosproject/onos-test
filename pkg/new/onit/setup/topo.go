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

// Topo is an interface for setting up topo nodes
type Topo interface {
	ServiceType
	concurrentSetup

	// Nodes sets the number of clusterTopo nodes to deploy
	Nodes(nodes int) Topo
}

var _ Topo = &clusterTopo{}

// clusterTopo is an implementation of the Topo interface
type clusterTopo struct {
	*clusterServiceType
	topo *cluster.Topo
}

func (s *clusterTopo) Nodes(nodes int) Topo {
	s.topo.SetNodes(nodes)
	return s
}

func (s *clusterTopo) create() error {
	return s.topo.Create()
}

// waitForStart waits for the onos-topo pods to complete startup
func (s *clusterTopo) waitForStart() error {
	return s.topo.AwaitReady()
}
