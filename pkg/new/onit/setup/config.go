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

// ConfigSetup is an interface for setting up config nodes
type ConfigSetup interface {
	Setup
	concurrentSetup
	Image(image string) ConfigSetup
	PullPolicy(pullPolicy corev1.PullPolicy) ConfigSetup
	Nodes(nodes int) ConfigSetup
}

var _ ConfigSetup = &configSetup{}

// configSetup is an implementation of the ConfigSetup interface
type configSetup struct {
	*testSetup
	image      string
	pullPolicy corev1.PullPolicy
	nodes      int
}

func (s *configSetup) Image(image string) ConfigSetup {
	s.image = image
	return s
}

func (s *configSetup) PullPolicy(pullPolicy corev1.PullPolicy) ConfigSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *configSetup) Nodes(nodes int) ConfigSetup {
	s.nodes = nodes
	return s
}

func (s *configSetup) create() error {
	return nil
}

func (s *configSetup) waitForStart() error {
	return nil
}
