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

import "github.com/onosproject/onos-test/pkg/new/onit/setup"

// Simulators provides the simulators environment
type Simulators interface {
	// Simulators returns a list of simulators in the environment
	Simulators() []Simulator

	// Get returns the environment for a simulator service by name
	Get(name string) Simulator

	// Add adds a new simulator to the environment
	Add(name string) setup.SimulatorSetup
}

var _ Simulators = &simulators{}

// simulators is an implementation of the Simulators interface
type simulators struct {
	*testEnv
}

func (e *simulators) Simulators() []Simulator {
	panic("implement me")
}

func (e *simulators) Get(name string) Simulator {
	return &simulator{
		service: &service{
			testEnv: e.testEnv,
		},
	}
}

func (e *simulators) Add(name string) setup.SimulatorSetup {
	return &simulatorSetup{
		name: name,
		serviceTypeSetup: &serviceTypeSetup{
			serviceSetup: &serviceSetup{
				testEnv: e.testEnv,
			},
		},
	}
}
