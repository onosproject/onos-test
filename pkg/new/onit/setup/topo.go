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
	corev1 "k8s.io/api/core/v1"
)

// TopoSetup is an interface for setting up topo nodes
type TopoSetup interface {
	Setup
	concurrentSetup
	Image(image string) TopoSetup
	PullPolicy(pullPolicy corev1.PullPolicy) TopoSetup
	Nodes(nodes int) TopoSetup
}

var _ TopoSetup = &topoSetup{}

// topoSetup is an implementation of the TopoSetup interface
type topoSetup struct {
	*testSetup
	image      string
	pullPolicy corev1.PullPolicy
	nodes      int
}

func (s *topoSetup) Image(image string) TopoSetup {
	s.image = image
	return s
}

func (s *topoSetup) PullPolicy(pullPolicy corev1.PullPolicy) TopoSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *topoSetup) Nodes(nodes int) TopoSetup {
	s.nodes = nodes
	return s
}

func (s *topoSetup) create() error {
	return nil
}

func (s *topoSetup) waitForStart() error {
	return nil
}
