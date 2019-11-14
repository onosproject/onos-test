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
	"github.com/openconfig/gnmi/client"
	gnmi "github.com/openconfig/gnmi/client/gnmi"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"time"
)

// SimulatorSetup is an interface for setting up a simulator
type SimulatorSetup interface {
	// Name sets the simulator name
	Name(name string) SimulatorSetup

	// Image sets the image to deploy
	Image(image string) SimulatorSetup

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) SimulatorSetup

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

func (s *clusterSimulatorSetup) Name(name string) SimulatorSetup {
	s.simulator.SetName(name)
	return s
}

func (s *clusterSimulatorSetup) Image(image string) SimulatorSetup {
	s.simulator.SetImage(image)
	return s
}

func (s *clusterSimulatorSetup) PullPolicy(pullPolicy corev1.PullPolicy) SimulatorSetup {
	s.simulator.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterSimulatorSetup) Add() (SimulatorEnv, error) {
	if err := s.simulator.Add(); err != nil {
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

func (e *clusterSimulatorEnv) Remove() error {
	return e.simulator.Remove()
}

func (e *clusterSimulatorEnv) RemoveOrDie() {
	if err := e.Remove(); err != nil {
		panic(err)
	}
}
