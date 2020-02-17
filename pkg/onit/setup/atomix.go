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

// AtomixSetup is an interface for setting up the Atomix controller
type AtomixSetup interface {
	// SetImage sets the Atomix controller image to deploy
	SetImage(image string) AtomixSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) AtomixSetup

	// SetReplicas sets the number of replicas to deploy
	SetReplicas(replicas int) AtomixSetup
}

var _ AtomixSetup = &clusterAtomixSetup{}

// clusterAtomixSetup is an implementation of the Atomix interface
type clusterAtomixSetup struct {
	atomix *cluster.Atomix
}

func (s *clusterAtomixSetup) SetImage(image string) AtomixSetup {
	s.atomix.SetImage(image)
	return s
}

func (s *clusterAtomixSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) AtomixSetup {
	s.atomix.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterAtomixSetup) SetReplicas(replicas int) AtomixSetup {
	s.atomix.SetReplicas(replicas)
	return s
}

func (s *clusterAtomixSetup) setup() error {
	return s.atomix.Setup()
}
