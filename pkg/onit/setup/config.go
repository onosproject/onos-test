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

// ConfigSetup is an interface for setting up config nodes
type ConfigSetup interface {
	// SetReplicas sets the number of replicas to deploy
	SetReplicas(replicas int) ConfigSetup

	// SetImage sets the onos-config image to deploy
	SetImage(image string) ConfigSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) ConfigSetup

	// SetCpuRequest sets the cpu request
	SetCPURequest(cpuRequest string) ConfigSetup

	// SetMemoryRequest sets memory request
	SetMemoryRequest(memoryRequest string) ConfigSetup

	// SetMemoryLimit sets memory limit
	SetMemoryLimit(memoryLimit string) ConfigSetup

	// SetCPULimit sets cpu limit
	SetCPULimit(cpuLimit string) ConfigSetup
}

var _ ConfigSetup = &clusterConfigSetup{}

// clusterConfigSetup is an implementation of the Config interface
type clusterConfigSetup struct {
	config *cluster.Config
}

func (s *clusterConfigSetup) SetReplicas(replicas int) ConfigSetup {
	s.config.SetReplicas(replicas)
	return s
}

func (s *clusterConfigSetup) SetImage(image string) ConfigSetup {
	s.config.SetImage(image)
	return s
}

func (s *clusterConfigSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) ConfigSetup {
	s.config.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterConfigSetup) SetCPURequest(cpuRequest string) ConfigSetup {
	s.config.SetCPURequest(cpuRequest)
	return s
}

func (s *clusterConfigSetup) SetMemoryRequest(memoryRequest string) ConfigSetup {
	s.config.SetMemoryRequest(memoryRequest)
	return s
}

func (s *clusterConfigSetup) SetMemoryLimit(memoryLimit string) ConfigSetup {
	s.config.SetMemoryLimit(memoryLimit)
	return s
}

func (s *clusterConfigSetup) SetCPULimit(cpuLimit string) ConfigSetup {
	s.config.SetCPULimit(cpuLimit)
	return s
}

func (s *clusterConfigSetup) setup() error {
	return s.config.Setup()
}
