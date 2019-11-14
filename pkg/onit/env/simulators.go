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
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"sync"
)

const simulatorQueueSize = 10

// SimulatorsSetup provides a setup configuration for multiple simulators
type SimulatorsSetup interface {
	// With adds a simulator setup
	With(setups ...SimulatorSetup) SimulatorsSetup

	// AddAll deploys the simulators in the cluster
	AddAll() ([]SimulatorEnv, error)

	// AddAllOrDie deploys the simulators and panics if the deployment fails
	AddAllOrDie() []SimulatorEnv
}

var _ SimulatorsSetup = &clusterSimulatorsSetup{}

// clusterSimulatorsSetup is an implementation of the SimulatorsSetup interface
type clusterSimulatorsSetup struct {
	simulators *cluster.Simulators
	setups     []SimulatorSetup
}

func (s *clusterSimulatorsSetup) With(setups ...SimulatorSetup) SimulatorsSetup {
	s.setups = append(s.setups, setups...)
	return s
}

func (s *clusterSimulatorsSetup) AddAll() ([]SimulatorEnv, error) {
	wg := &sync.WaitGroup{}
	wg.Add(len(s.setups))

	envCh := make(chan SimulatorEnv, len(s.setups))
	errCh := make(chan error)

	setupQueue := make(chan SimulatorSetup, simulatorQueueSize)
	for i := 0; i < simulatorQueueSize; i++ {
		go func() {
			for setup := range setupQueue {
				if simulator, err := setup.Add(); err != nil {
					errCh <- err
				} else {
					envCh <- simulator
				}
				wg.Done()
			}
		}()
	}

	go func() {
		for _, setup := range s.setups {
			setupQueue <- setup
		}
		close(setupQueue)
	}()

	go func() {
		wg.Wait()
		close(envCh)
		close(errCh)
	}()

	for err := range errCh {
		return nil, err
	}

	simulators := make([]SimulatorEnv, 0, len(s.setups))
	for simulator := range envCh {
		simulators = append(simulators, simulator)
	}
	return simulators, nil
}

func (s *clusterSimulatorsSetup) AddAllOrDie() []SimulatorEnv {
	if simulators, err := s.AddAll(); err != nil {
		panic(err)
	} else {
		return simulators
	}
}

// SimulatorsEnv provides the simulators environment
type SimulatorsEnv interface {
	// List returns a list of simulators in the environment
	List() []SimulatorEnv

	// Get returns the environment for a simulator service by name
	Get(name string) SimulatorEnv

	// New adds a new simulator to the environment
	New() SimulatorSetup
}

var _ SimulatorsEnv = &clusterSimulatorsEnv{}

// clusterSimulatorsEnv is an implementation of the Simulators interface
type clusterSimulatorsEnv struct {
	simulators *cluster.Simulators
}

func (e *clusterSimulatorsEnv) List() []SimulatorEnv {
	clusterSimulators := e.simulators.List()
	simulators := make([]SimulatorEnv, len(clusterSimulators))
	for i, simulator := range clusterSimulators {
		simulators[i] = e.Get(simulator.Name())
	}
	return simulators
}

func (e *clusterSimulatorsEnv) Get(name string) SimulatorEnv {
	simulator := e.simulators.Get(name)
	return &clusterSimulatorEnv{
		clusterNodeEnv: &clusterNodeEnv{
			node: simulator.Node,
		},
		simulator: simulator,
	}
}

func (e *clusterSimulatorsEnv) New() SimulatorSetup {
	return &clusterSimulatorSetup{
		simulator: e.simulators.New(),
	}
}
