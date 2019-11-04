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
)

func newSimulatorSetup(name string, testEnv *testEnv) setup.SimulatorSetup {
	setup := &simulatorSetup{
		serviceSetup: &serviceSetup{
			testEnv: testEnv,
		},
		name: name,
	}
	setup.serviceSetup.setup = setup
	return setup
}

// Simulator provides the environment for a single simulator
type Simulator interface {
	Service
}

var _ Simulator = &simulator{}

// simulator is an implementation of the Simulator interface
type simulator struct {
	*service
}

var _ setup.SimulatorSetup = &simulatorSetup{}

// simulatorSetup is an implementation of the SimulatorSetup interface
type simulatorSetup struct {
	*serviceSetup
	name string
}

func (s *simulatorSetup) Using() setup.ServiceSetup {
	return s
}

func (s *simulatorSetup) Setup() error {
	return nil
}
