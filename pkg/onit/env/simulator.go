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
	"context"
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-topo/api/device"
	"github.com/openconfig/gnmi/client"
	gnmi "github.com/openconfig/gnmi/client/gnmi"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"time"
)

// DevicePredicate is a function for evaluating the state of a device
type DevicePredicate func(*device.Device) bool

// SimulatorSetup is an interface for setting up a simulator
type SimulatorSetup interface {
	// SetName sets the simulator name
	SetName(name string) SimulatorSetup

	// SetImage sets the image to deploy
	SetImage(image string) SimulatorSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) SimulatorSetup

	// SetDeviceType sets the device type
	SetDeviceType(deviceType string) SimulatorSetup

	// SetDeviceVersion sets the device version
	SetDeviceVersion(version string) SimulatorSetup

	// SetDeviceTimeout sets the device timeout
	SetDeviceTimeout(timeout time.Duration) SimulatorSetup

	// Add deploys the simulator in the cluster
	Add() (SimulatorEnv, error)

	// AddOrDie deploys the simulator and panics if the deployment fails
	AddOrDie() SimulatorEnv
}

var _ SimulatorSetup = &clusterSimulatorSetup{}

// clusterSimulatorSetup is an implementation of the SimulatorSetup interface
type clusterSimulatorSetup struct {
	simulator *cluster.Simulator
}

func (s *clusterSimulatorSetup) SetName(name string) SimulatorSetup {
	s.simulator.SetName(name)
	return s
}

func (s *clusterSimulatorSetup) SetImage(image string) SimulatorSetup {
	s.simulator.SetImage(image)
	return s
}

func (s *clusterSimulatorSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) SimulatorSetup {
	s.simulator.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterSimulatorSetup) SetDeviceType(deviceType string) SimulatorSetup {
	s.simulator.SetDeviceType(deviceType)
	return s
}

func (s *clusterSimulatorSetup) SetDeviceVersion(version string) SimulatorSetup {
	s.simulator.SetDeviceVersion(version)
	return s
}

func (s *clusterSimulatorSetup) SetDeviceTimeout(timeout time.Duration) SimulatorSetup {
	s.simulator.SetDeviceTimeout(timeout)
	return s
}

func (s *clusterSimulatorSetup) Add() (SimulatorEnv, error) {
	if err := s.simulator.Setup(); err != nil {
		return nil, err
	}
	return &clusterSimulatorEnv{
		clusterNodeEnv: &clusterNodeEnv{
			node: s.simulator.Node,
		},
		simulator: s.simulator,
	}, nil
}

func (s *clusterSimulatorSetup) AddOrDie() SimulatorEnv {
	network, err := s.Add()
	if err != nil {
		panic(err)
	}
	return network
}

// SimulatorEnv provides the environment for a single simulator
type SimulatorEnv interface {
	NodeEnv

	// Destination returns the gNMI client destination
	Destination() client.Destination

	// NewGNMIClient returns the gNMI client
	NewGNMIClient() (*gnmi.Client, error)

	// Await waits for the simulator device state to match the given predicate
	Await(predicate DevicePredicate, timeout time.Duration) error

	// Remove removes the simulator
	Remove() error

	// RemoveOrDie removes the simulator and panics if the remove fails
	RemoveOrDie()
}

var _ SimulatorEnv = &clusterSimulatorEnv{}

// clusterSimulatorEnv is an implementation of the Simulator interface
type clusterSimulatorEnv struct {
	*clusterNodeEnv
	simulator *cluster.Simulator
}

func (e *clusterSimulatorEnv) Connect() (*grpc.ClientConn, error) {
	return grpc.Dial(e.Address(), grpc.WithInsecure())
}

func (e *clusterSimulatorEnv) Destination() client.Destination {
	return client.Destination{
		Addrs:   []string{e.Address()},
		Target:  e.Name(),
		TLS:     e.Credentials(),
		Timeout: 10 * time.Second,
	}
}

func (e *clusterSimulatorEnv) NewGNMIClient() (*gnmi.Client, error) {
	conn, err := e.Connect()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := gnmi.NewFromConn(ctx, conn, e.Destination())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (e *clusterSimulatorEnv) Await(predicate DevicePredicate, timeout time.Duration) error {
	return e.simulator.AwaitDevicePredicate(predicate, timeout)
}

func (e *clusterSimulatorEnv) Remove() error {
	return e.simulator.TearDown()
}

func (e *clusterSimulatorEnv) RemoveOrDie() {
	if err := e.Remove(); err != nil {
		panic(err)
	}
}
