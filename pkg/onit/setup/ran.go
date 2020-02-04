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

// RanSetup is an interface for setting up ran nodes
type RanSetup interface {

	// SetEnabled enables the Ran subsystem
	SetEnabled() RanSetup

	// SetReplicas sets the number of ran replicas to deploy
	SetReplicas(replicas int) RanSetup

	// SetImage sets the onos-ran image to deploy
	SetImage(image string) RanSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) RanSetup

	// SetCpuRequest sets the cpu request
	SetCPURequest(cpuRequest string) RanSetup

	// SetMemoryRequest sets memory request
	SetMemoryRequest(memoryRequest string) RanSetup

	// SetMemoryLimit sets memory limit
	SetMemoryLimit(memoryLimit string) RanSetup

	// SetCPULimit sets cpu limit
	SetCPULimit(cpuLimit string) RanSetup
}

var _ RanSetup = &clusterRanSetup{}

// clusterRanSetup is an implementation of the Ran interface
type clusterRanSetup struct {
	ran *cluster.Ran
}

func (s *clusterRanSetup) SetEnabled() RanSetup {
	s.ran.SetEnabled(true)
	return s
}

func (s *clusterRanSetup) SetReplicas(replicas int) RanSetup {
	s.ran.SetReplicas(replicas)
	return s
}

func (s *clusterRanSetup) SetImage(image string) RanSetup {
	s.ran.SetImage(image)
	return s
}

func (s *clusterRanSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) RanSetup {
	s.ran.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterRanSetup) SetCPURequest(cpuRequest string) RanSetup {
	s.ran.SetCPURequest(cpuRequest)
	return s
}

func (s *clusterRanSetup) SetMemoryRequest(memoryRequest string) RanSetup {
	s.ran.SetMemoryRequest(memoryRequest)
	return s
}

func (s *clusterRanSetup) SetMemoryLimit(memoryLimit string) RanSetup {
	s.ran.SetMemoryLimit(memoryLimit)
	return s
}

func (s *clusterRanSetup) SetCPULimit(cpuLimit string) RanSetup {
	s.ran.SetCPULimit(cpuLimit)
	return s
}

func (s *clusterRanSetup) setup() error {
	return s.ran.Setup()
}
