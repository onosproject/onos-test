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
	corev1 "k8s.io/api/core/v1"
)

// Atomix is an interface for setting up the Atomix controller
type Atomix interface {
	// Image sets the Atomix controller image to deploy
	Image(image string) Atomix

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) Atomix
}

var _ Atomix = &clusterAtomix{}

// clusterAtomix is an implementation of the Atomix interface
type clusterAtomix struct {
	atomix *cluster.Atomix
}

func (s *clusterAtomix) Image(image string) Atomix {
	s.atomix.SetImage(image)
	return s
}

func (s *clusterAtomix) PullPolicy(pullPolicy corev1.PullPolicy) Atomix {
	s.atomix.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterAtomix) setup() error {
	return s.atomix.Setup()
}
