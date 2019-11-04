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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Simulators provides the simulators environment
type Simulators interface {
	// List returns a list of simulators in the environment
	List() []Simulator

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

func (e *simulators) List() []Simulator {
	pods, err := e.kubeClient.CoreV1().Pods(e.namespace).List(metav1.ListOptions{
		LabelSelector: "type=simulator",
	})
	if err != nil {
		panic(err)
	}

	simulators := make([]Simulator, len(pods.Items))
	for i, pod := range pods.Items {
		simulators[i] = e.Get(pod.Name)
	}
	return simulators
}

func (e *simulators) Get(name string) Simulator {
	return &simulator{
		service: newService(name, "simulator", e.testEnv),
	}
}

func (e *simulators) Add(name string) setup.SimulatorSetup {
	return newSimulatorSetup(name, e.testEnv)
}
