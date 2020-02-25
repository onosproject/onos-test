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

// RANSetup is an interface for setting up ran nodes
type RANSetup interface {

	// SetEnabled enables the Ran subsystem
	SetEnabled() RANSetup

	// SetReplicas sets the number of ran replicas to deploy
	SetReplicas(replicas int) RANSetup

	// SetImage sets the onos-ric image to deploy
	SetImage(image string) RANSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) RANSetup

	// SetCpuRequest sets the cpu request
	SetCPURequest(cpuRequest string) RANSetup

	// SetMemoryRequest sets memory request
	SetMemoryRequest(memoryRequest string) RANSetup

	// SetMemoryLimit sets memory limit
	SetMemoryLimit(memoryLimit string) RANSetup

	// SetCPULimit sets cpu limit
	SetCPULimit(cpuLimit string) RANSetup
}

var _ RANSetup = &clusterRANSetup{}

// clusterRANSetup is an implementation of the Ran interface
type clusterRANSetup struct {
	ran *cluster.RAN
}

func (s *clusterRANSetup) SetEnabled() RANSetup {
	s.ran.SetEnabled(true)
	return s
}

func (s *clusterRANSetup) SetReplicas(replicas int) RANSetup {
	s.ran.SetReplicas(replicas)
	return s
}

func (s *clusterRANSetup) SetImage(image string) RANSetup {
	s.ran.SetImage(image)
	return s
}

func (s *clusterRANSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) RANSetup {
	s.ran.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterRANSetup) SetCPURequest(cpuRequest string) RANSetup {
	s.ran.SetCPURequest(cpuRequest)
	return s
}

func (s *clusterRANSetup) SetMemoryRequest(memoryRequest string) RANSetup {
	s.ran.SetMemoryRequest(memoryRequest)
	return s
}

func (s *clusterRANSetup) SetMemoryLimit(memoryLimit string) RANSetup {
	s.ran.SetMemoryLimit(memoryLimit)
	return s
}

func (s *clusterRANSetup) SetCPULimit(cpuLimit string) RANSetup {
	s.ran.SetCPULimit(cpuLimit)
	return s
}

func (s *clusterRANSetup) setup() error {
	return s.ran.Setup()
}
