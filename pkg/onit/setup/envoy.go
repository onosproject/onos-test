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

// EnvoySetup is an interface for setting up envoy nodes
type EnvoySetup interface {

	// SetEnabled enables the Envoy
	SetEnabled() EnvoySetup

	// SetReplicas sets the number of envoy replicas to deploy
	SetReplicas(replicas int) EnvoySetup

	// SetImage sets the envoy image to deploy
	SetImage(image string) EnvoySetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) EnvoySetup

	// SetConfigMaps sets the config maps
	SetConfigMaps(map[string]string) EnvoySetup
}

var _ EnvoySetup = &clusterEnvoySetup{}

// clusterGuiSetup is an implementation of the Gui interface
type clusterEnvoySetup struct {
	envoy *cluster.Envoy
}

func (s *clusterEnvoySetup) SetEnabled() EnvoySetup {
	s.envoy.SetEnabled(true)
	return s
}

func (s *clusterEnvoySetup) SetReplicas(replicas int) EnvoySetup {
	s.envoy.SetReplicas(replicas)
	return s
}

func (s *clusterEnvoySetup) SetImage(image string) EnvoySetup {
	s.envoy.SetImage(image)
	return s
}

func (s *clusterEnvoySetup) SetPullPolicy(pullPolicy corev1.PullPolicy) EnvoySetup {
	s.envoy.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterEnvoySetup) SetConfigMaps(configMaps map[string]string) EnvoySetup {
	s.envoy.SetConfigMaps(configMaps)
	return s
}

func (s *clusterEnvoySetup) setup() error {
	if s.envoy.Enabled() {
		return s.envoy.Setup()
	}
	return nil
}
