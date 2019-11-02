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

package env

import (
	"github.com/onosproject/onos-test/pkg/new/onit/setup"
	corev1 "k8s.io/api/core/v1"
)

// SimulatorEnv provides the environment for a single simulator
type SimulatorEnv interface {
	ServiceEnv
}

var _ SimulatorEnv = &simulatorEnv{}

// simulatorEnv is an implementation of the SimulatorEnv interface
type simulatorEnv struct {
	*serviceEnv
}

var _ setup.SimulatorSetup = &simulatorSetup{}

// simulatorSetup is an implementation of the SimulatorSetup interface
type simulatorSetup struct {
	*testEnv
	name       string
	image      string
	pullPolicy corev1.PullPolicy
}

func (s *simulatorSetup) Name(name string) setup.SimulatorSetup {
	s.name = name
	return s
}

func (s *simulatorSetup) Image(image string) setup.SimulatorSetup {
	s.image = image
	return s
}

func (s *simulatorSetup) PullPolicy(pullPolicy corev1.PullPolicy) setup.SimulatorSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *simulatorSetup) Setup() error {
	return nil
}
