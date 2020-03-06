// Copyright 2020-present Open Networking Foundation.
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

// RICSetup is an interface for setting up ran nodes
type RICSetup interface {

	// SetEnabled enables the Ran subsystem
	SetEnabled() RICSetup

	// SetReplicas sets the number of ran replicas to deploy
	SetReplicas(replicas int) RICSetup

	// SetImage sets the onos-ric image to deploy
	SetImage(image string) RICSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) RICSetup

	// SetCpuRequest sets the cpu request
	SetCPURequest(cpuRequest string) RICSetup

	// SetMemoryRequest sets memory request
	SetMemoryRequest(memoryRequest string) RICSetup

	// SetMemoryLimit sets memory limit
	SetMemoryLimit(memoryLimit string) RICSetup

	// SetCPULimit sets cpu limit
	SetCPULimit(cpuLimit string) RICSetup
}

var _ RICSetup = &clusterRICSetup{}

// clusterRICSetup is an implementation of the Ran interface
type clusterRICSetup struct {
	ran *cluster.RIC
}

func (s *clusterRICSetup) SetEnabled() RICSetup {
	s.ran.SetEnabled(true)
	return s
}

func (s *clusterRICSetup) SetReplicas(replicas int) RICSetup {
	s.ran.SetReplicas(replicas)
	return s
}

func (s *clusterRICSetup) SetImage(image string) RICSetup {
	s.ran.SetImage(image)
	return s
}

func (s *clusterRICSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) RICSetup {
	s.ran.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterRICSetup) SetCPURequest(cpuRequest string) RICSetup {
	s.ran.SetCPURequest(cpuRequest)
	return s
}

func (s *clusterRICSetup) SetMemoryRequest(memoryRequest string) RICSetup {
	s.ran.SetMemoryRequest(memoryRequest)
	return s
}

func (s *clusterRICSetup) SetMemoryLimit(memoryLimit string) RICSetup {
	s.ran.SetMemoryLimit(memoryLimit)
	return s
}

func (s *clusterRICSetup) SetCPULimit(cpuLimit string) RICSetup {
	s.ran.SetCPULimit(cpuLimit)
	return s
}

func (s *clusterRICSetup) setup() error {
	return s.ran.Setup()
}
