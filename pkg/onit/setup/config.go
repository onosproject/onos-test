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

// Config is an interface for setting up config nodes
type Config interface {
	// Nodes sets the number of nodes to deploy
	Nodes(nodes int) Config

	// Image sets the onos-config image to deploy
	Image(image string) Config

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) Config
}

var _ Config = &clusterConfig{}

// clusterConfig is an implementation of the Config interface
type clusterConfig struct {
	config *cluster.Config
}

func (s *clusterConfig) Nodes(nodes int) Config {
	s.config.SetReplicas(nodes)
	return s
}

func (s *clusterConfig) Image(image string) Config {
	s.config.SetImage(image)
	return s
}

func (s *clusterConfig) PullPolicy(pullPolicy corev1.PullPolicy) Config {
	s.config.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterConfig) create() error {
	return s.config.Create()
}

func (s *clusterConfig) waitForStart() error {
	return s.config.AwaitReady()
}
