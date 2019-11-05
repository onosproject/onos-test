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

package cluster

import "github.com/onosproject/onos-test/pkg/new/util/random"

func newSimulators(client *client) *Simulators {
	return &Simulators{
		client: client,
	}
}

// Simulators provides methods for adding and modifying simulators
type Simulators struct {
	*client
}

// New returns a new simulator
func (s *Simulators) New() *Simulator {
	return newSimulator(random.NewPetName(2), s.client)
}

// Get gets a simulator by name
func (s *Simulators) Get(name string) *Simulator {
	return newSimulator(name, s.client)
}

// List lists the simulators in the cluster
func (s *Simulators) List() []*Simulator {
	names := s.listServices(simulatorType)
	simulators := make([]*Simulator, len(names))
	for i, name := range names {
		simulators[i] = s.Get(name)
	}
	return simulators
}
