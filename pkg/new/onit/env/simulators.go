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

// SimulatorsEnv provides the simulators environment
type SimulatorsEnv interface {
	Simulator(name string) SimulatorEnv
	Add() setup.SimulatorSetup
}

var _ SimulatorsEnv = &simulatorsEnv{}

// simulatorsEnv is an implementation of the SimulatorsEnv interface
type simulatorsEnv struct {
	*testEnv
}

func (e *simulatorsEnv) Simulator(name string) SimulatorEnv {
	return &simulatorEnv{
		serviceEnv: &serviceEnv{
			testEnv: e.testEnv,
		},
	}
}

func (e *simulatorsEnv) Add() setup.SimulatorSetup {
	return &simulatorSetup{
		testEnv: e.testEnv,
	}
}
