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

// AddNetwork add a stratum network to the cluster
func (t *TestSetup) AddNetwork() {

	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	// Create the network configuration
	config := &k8s.NetworkConfig{
		Config: t.configName,
	}

	// Update number of devices in the network configuration
	k8s.ParseMininetOptions(config)

	if err != nil {
		exitError(err)
	}

	// Add the network to the cluster
	if status := cluster.AddNetwork(t.networkName, config); status.Failed() {
		exitStatus(status)
	} else {
		fmt.Println(t.networkName)
	}

}

// RemoveNetwork remove a network from the cluster
func (t *TestSetup) RemoveNetwork() {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	networks, err := cluster.GetNetworks()
	if err != nil {
		exitError(err)
	}
	if !Contains(networks, t.networkName) {
		exitError(errors.New("The given network name does not exist"))
	}

	// Remove the network from the cluster
	if status := cluster.RemoveNetwork(t.networkName); status.Failed() {
		exitStatus(status)
	}
}

// GetNetworks returns the list of networks in the cluster
func (t *TestSetup) GetNetworks() ([]string, error) {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}
	return cluster.GetNetworks()
}
