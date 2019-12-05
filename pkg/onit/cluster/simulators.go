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

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/util/random"
)

func newSimulators(cluster *Cluster) *Simulators {
	return &Simulators{
		client:  cluster.client,
		cluster: cluster,
	}
}

// Simulators provides methods for adding and modifying simulators
type Simulators struct {
	*client
	cluster *Cluster
}

// New returns a new simulator
func (s *Simulators) New() *Simulator {
	return newSimulator(s.cluster, fmt.Sprintf("devicesim-%s", random.NewPetName(2)))
}

// Get gets a simulator by name
func (s *Simulators) Get(name string) *Simulator {
	return newSimulator(s.cluster, name)
}

// List lists the simulators in the cluster
func (s *Simulators) List() []*Simulator {
	names := s.listServices(getLabels(simulatorType))
	simulators := make([]*Simulator, len(names))
	for i, name := range names {
		simulators[i] = s.Get(name)
	}
	return simulators
}
