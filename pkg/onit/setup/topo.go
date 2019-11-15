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
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	corev1 "k8s.io/api/core/v1"
)

// TopoSetup is an interface for setting up topo nodes
type TopoSetup interface {
	// SetNodes sets the number of clusterTopo nodes to deploy
	SetNodes(nodes int) TopoSetup

	// SetImage sets the onos-topo image to deploy
	SetImage(image string) TopoSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) TopoSetup
}

var _ TopoSetup = &clusterTopoSetup{}

// clusterTopoSetup is an implementation of the Topo interface
type clusterTopoSetup struct {
	topo *cluster.Topo
}

func (s *clusterTopoSetup) SetNodes(nodes int) TopoSetup {
	s.topo.SetReplicas(nodes)
	return s
}

func (s *clusterTopoSetup) SetImage(image string) TopoSetup {
	s.topo.SetImage(image)
	return s
}

func (s *clusterTopoSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) TopoSetup {
	s.topo.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterTopoSetup) setup() error {
	return s.topo.Setup()
}
