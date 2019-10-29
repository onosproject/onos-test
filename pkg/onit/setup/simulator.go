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
	"errors"
	"fmt"

	"github.com/onosproject/onos-test/pkg/onit/k8s"
)

// GetSimulators get list of simulators in the current cluster
func (t *TestSetup) GetSimulators() ([]string, error) {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}
	// Get the list of simulators and output
	return cluster.GetSimulators()
}

// AddSimulator add a simulator to the cluster
func (t *TestSetup) AddSimulator() {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	// Create the simulator configuration
	config := &k8s.SimulatorConfig{
		Config: t.configName,
	}

	// Add the simulator to the cluster
	if status := cluster.AddSimulator(t.simulatorName, config); status.Failed() {
		exitStatus(status)
	} else {
		fmt.Println(t.simulatorName)
	}
}

// RemoveSimulator remove a simulator based on its given name
func (t *TestSetup) RemoveSimulator() {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	simulators, err := cluster.GetSimulators()
	if err != nil {
		exitError(err)
	}

	if !Contains(simulators, t.simulatorName) {
		exitError(errors.New("The given simulator name does not exist"))
	}

	// Remove the simulator from the cluster
	if status := cluster.RemoveSimulator(t.simulatorName); status.Failed() {
		exitStatus(status)
	}
}
