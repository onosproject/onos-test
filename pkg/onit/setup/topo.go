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

// Topo is an interface for setting up topo nodes
type Topo interface {
	// Nodes sets the number of clusterTopo nodes to deploy
	Nodes(nodes int) Topo

	// Image sets the onos-topo image to deploy
	Image(image string) Topo

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) Topo
}

var _ Topo = &clusterTopo{}

// clusterTopo is an implementation of the Topo interface
type clusterTopo struct {
	topo *cluster.Topo
}

func (s *clusterTopo) Nodes(nodes int) Topo {
	s.topo.SetReplicas(nodes)
	return s
}

func (s *clusterTopo) Image(image string) Topo {
	s.topo.SetImage(image)
	return s
}

func (s *clusterTopo) PullPolicy(pullPolicy corev1.PullPolicy) Topo {
	s.topo.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterTopo) create() error {
	return s.topo.Create()
}

// waitForStart waits for the onos-topo pods to complete startup
func (s *clusterTopo) waitForStart() error {
	return s.topo.AwaitReady()
}
