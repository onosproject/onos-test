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
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
	corev1 "k8s.io/api/core/v1"
)

// NetworkSetup is an interface for deploying up a network
type NetworkSetup interface {
	// Name sets the network name
	Name(name string) NetworkSetup

	// Single creates a single node topology
	Single() NetworkSetup

	// Linear creates a linear topology with the given number of devices
	Linear(devices int) NetworkSetup

	// Topo creates a custom topology
	Topo(topo string, devices int) NetworkSetup

	// Image sets the image to deploy
	Image(image string) NetworkSetup

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) NetworkSetup

	// Add adds the network to the cluster
	Add() (Network, error)

	// AddOrDie adds the network and panics if the deployment fails
	AddOrDie() Network
}

var _ NetworkSetup = &clusterNetworkSetup{}

// clusterNetworkSetup is an implementation of the NetworkSetup interface
type clusterNetworkSetup struct {
	network *cluster.Network
}

func (s *clusterNetworkSetup) Name(name string) NetworkSetup {
	s.network.SetName(name)
	return s
}

func (s *clusterNetworkSetup) Image(image string) NetworkSetup {
	s.network.SetImage(image)
	return s
}

func (s *clusterNetworkSetup) PullPolicy(pullPolicy corev1.PullPolicy) NetworkSetup {
	s.network.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterNetworkSetup) Single() NetworkSetup {
	s.network.SetSingle()
	return s
}

func (s *clusterNetworkSetup) Linear(devices int) NetworkSetup {
	s.network.SetLinear(devices)
	return s
}

func (s *clusterNetworkSetup) Topo(topo string, devices int) NetworkSetup {
	s.network.SetTopo(topo, devices)
	return s
}

func (s *clusterNetworkSetup) Add() (Network, error) {
	if err := s.network.Add(); err != nil {
		return nil, err
	}
	return &clusterNetwork{
		clusterNode: &clusterNode{
			node: s.network.Node,
		},
		network: s.network,
	}, nil
}

func (s *clusterNetworkSetup) AddOrDie() Network {
	network, err := s.Add()
	if err != nil {
		panic(err)
	}
	return network
}

// Network provides the environment for a network node
type Network interface {
	Node

	// Devices returns a list of devices in the network
	Devices() []Node

	// Remove removes the network
	Remove() error

	// RemoveOrDie removes the network and panics if the remove fails
	RemoveOrDie()
}

var _ Network = &clusterNetwork{}

// clusterNetwork is an implementation of the Network interface
type clusterNetwork struct {
	*clusterNode
	network *cluster.Network
}

func (e *clusterNetwork) Devices() []Node {
	clusterDevices := e.network.Devices()
	devices := make([]Node, len(clusterDevices))
	for i, node := range clusterDevices {
		devices[i] = &clusterNode{
			node: node,
		}
	}
	return devices
}

func (e *clusterNetwork) Remove() error {
	return e.network.Remove()
}

func (e *clusterNetwork) RemoveOrDie() {
	if err := e.Remove(); err != nil {
		panic(err)
	}
}
