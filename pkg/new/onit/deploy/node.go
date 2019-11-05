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
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
	corev1 "k8s.io/api/core/v1"
)

// NodeType provides methods for setting up a single-node service
type NodeType interface {
	// Using returns the service setup
	Using() Node
}

// Node provides methods for setting up a single-node service
type Node interface {
	Deploy

	// Image sets the simulator image to deploy
	Image(image string) Node

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) Node
}

// clusterNodeType is an implementation of the NodeType interface
type clusterNodeType struct {
	*clusterNode
}

func (s *clusterNodeType) Using() Node {
	return s.clusterNode
}

var _ Node = &clusterNode{}

// clusterNode is an implementation of the Node interface
type clusterNode struct {
	node   *cluster.Node
	deploy Deploy
}

func (s *clusterNode) Image(image string) Node {
	s.node.SetImage(image)
	return s
}

func (s *clusterNode) PullPolicy(pullPolicy corev1.PullPolicy) Node {
	s.node.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterNode) Deploy() error {
	return s.deploy.Deploy()
}

func (s *clusterNode) DeployOrDie() {
	if err := s.Deploy(); err != nil {
		panic(err)
	}
}
