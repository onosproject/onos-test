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
	corev1 "k8s.io/api/core/v1"
	"time"
)

// NetworkSetup is an interface for deploying up a network
type NetworkSetup interface {
	// SetName sets the network name
	SetName(name string) NetworkSetup

	// SetSingle creates a single node topology
	SetSingle() NetworkSetup

	// SetLinear creates a linear topology with the given number of devices
	SetLinear(devices int) NetworkSetup

	// SetCustom creates a custom topology
	SetCustom(topo string, devices int) NetworkSetup

	// SetImage sets the image to deploy
	SetImage(image string) NetworkSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) NetworkSetup

	// SetDeviceType sets the device type
	SetDeviceType(deviceType string) NetworkSetup

	// SetDeviceVersion sets the device version
	SetDeviceVersion(version string) NetworkSetup

	// SetDeviceTimeout sets the device timeout
	SetDeviceTimeout(timeout time.Duration) NetworkSetup

	// Add adds the network to the cluster
	Add() (NetworkEnv, error)

	// AddOrDie adds the network and panics if the deployment fails
	AddOrDie() NetworkEnv
}

var _ NetworkSetup = &clusterNetworkSetup{}

// clusterNetworkSetup is an implementation of the NetworkSetup interface
type clusterNetworkSetup struct {
	network *cluster.Network
}

func (s *clusterNetworkSetup) SetName(name string) NetworkSetup {
	s.network.SetName(name)
	return s
}

func (s *clusterNetworkSetup) SetImage(image string) NetworkSetup {
	s.network.SetImage(image)
	return s
}

func (s *clusterNetworkSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) NetworkSetup {
	s.network.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterNetworkSetup) SetDeviceType(deviceType string) NetworkSetup {
	s.network.SetDeviceType(deviceType)
	return s
}

func (s *clusterNetworkSetup) SetDeviceVersion(version string) NetworkSetup {
	s.network.SetDeviceVersion(version)
	return s
}

func (s *clusterNetworkSetup) SetDeviceTimeout(timeout time.Duration) NetworkSetup {
	s.network.SetDeviceTimeout(timeout)
	return s
}

func (s *clusterNetworkSetup) SetSingle() NetworkSetup {
	s.network.SetSingle()
	return s
}

func (s *clusterNetworkSetup) SetLinear(devices int) NetworkSetup {
	s.network.SetLinear(devices)
	return s
}

func (s *clusterNetworkSetup) SetCustom(topo string, devices int) NetworkSetup {
	s.network.SetTopo(topo, devices)
	return s
}

func (s *clusterNetworkSetup) Add() (NetworkEnv, error) {
	if err := s.network.Setup(); err != nil {
		return nil, err
	}
	return &clusterNetworkEnv{
		clusterNodeEnv: &clusterNodeEnv{
			node: s.network.Node,
		},
		network: s.network,
	}, nil
}

func (s *clusterNetworkSetup) AddOrDie() NetworkEnv {
	network, err := s.Add()
	if err != nil {
		panic(err)
	}
	return network
}

// NetworkEnv provides the environment for a network node
type NetworkEnv interface {
	NodeEnv

	// Devices returns a list of devices in the network
	Devices() []NodeEnv

	// Remove removes the network
	Remove() error

	// RemoveOrDie removes the network and panics if the remove fails
	RemoveOrDie()
}

var _ NetworkEnv = &clusterNetworkEnv{}

// clusterNetworkEnv is an implementation of the Network interface
type clusterNetworkEnv struct {
	*clusterNodeEnv
	network *cluster.Network
}

func (e *clusterNetworkEnv) Devices() []NodeEnv {
	clusterDevices, err := e.network.Devices()
	if err != nil {
		panic(err)
	}

	devices := make([]NodeEnv, len(clusterDevices))
	for i, node := range clusterDevices {
		devices[i] = &clusterNodeEnv{
			node: node,
		}
	}
	return devices
}

func (e *clusterNetworkEnv) Remove() error {
	return e.network.TearDown()
}

func (e *clusterNetworkEnv) RemoveOrDie() {
	if err := e.Remove(); err != nil {
		panic(err)
	}
}
