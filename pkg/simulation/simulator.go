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

package simulation

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/model"
	"github.com/onosproject/onos-test/pkg/registry"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"google.golang.org/grpc"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// newSimulator returns a new simulation worker
func newSimulator(config *Config) (*Simulator, error) {
	kubeAPI, err := kube.GetAPI(config.ID)
	if err != nil {
		return nil, err
	}
	return &Simulator{
		client:      kubeAPI.Client(),
		simulations: make(map[string]*Simulation),
	}, nil
}

// Simulator runs a simulation job
type Simulator struct {
	client      client.Client
	simulations map[string]*Simulation
}

// Run runs a simulation
func (s *Simulator) Run() error {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	RegisterSimulatorServiceServer(server, s)
	return server.Serve(lis)
}

func (s *Simulator) getSimulation(name string, args map[string]string) (*Simulation, error) {
	if simulation, ok := s.simulations[name]; ok {
		return simulation, nil
	}
	suite := registry.GetSimulationSuite(name)
	if suite != nil {
		simulation := newSimulation(name, getSimulatorID(), suite, args)
		s.simulations[name] = simulation
		return simulation, nil
	}
	return nil, fmt.Errorf("unknown simulation %s", name)
}

// SetupSimulation sets up a simulation suite
func (s *Simulator) SetupSimulation(ctx context.Context, request *SimulationLifecycleRequest) (*SimulationLifecycleResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulatorID()), "SetupSimulation %s", request.Simulation)
	step.Start()

	simulation, err := s.getSimulation(request.Simulation, request.Args)
	if err != nil {
		step.Fail(err)
		return nil, err
	}
	simulation.setupSimulation()
	step.Complete()
	return &SimulationLifecycleResponse{}, nil
}

// TearDownSimulation tears down a simulation suite
func (s *Simulator) TearDownSimulation(ctx context.Context, request *SimulationLifecycleRequest) (*SimulationLifecycleResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulatorID()), "TearDownSimulation %s", request.Simulation)
	step.Start()

	simulation, err := s.getSimulation(request.Simulation, request.Args)
	if err != nil {
		step.Fail(err)
		return nil, err
	}
	simulation.teardownSimulation()
	step.Complete()
	return &SimulationLifecycleResponse{}, nil
}

func (s *Simulator) SetupSimulator(ctx context.Context, request *SimulationLifecycleRequest) (*SimulationLifecycleResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulatorID()), "SetupSimulator %s", request.Simulation)
	step.Start()

	simulation, err := s.getSimulation(request.Simulation, request.Args)
	if err != nil {
		step.Fail(err)
		return nil, err
	}
	simulation.setupSimulator()
	step.Complete()
	return &SimulationLifecycleResponse{}, nil
}

func (s *Simulator) TearDownSimulator(ctx context.Context, request *SimulationLifecycleRequest) (*SimulationLifecycleResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulatorID()), "TearDownSimulator %s", request.Simulation)
	step.Start()

	simulation, err := s.getSimulation(request.Simulation, request.Args)
	if err != nil {
		step.Fail(err)
		return nil, err
	}
	simulation.teardownSimulator()
	step.Complete()
	return &SimulationLifecycleResponse{}, nil
}

func (s *Simulator) Simulate(request *SimulateRequest, stream SimulatorService_SimulateServer) error {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulatorID()), "Simulate %s/%s", request.Simulation, request.Method)
	step.Start()

	simulation, ok := s.simulations[request.Simulation]
	if !ok {
		err := fmt.Errorf("unknown simulation %s", request.Simulation)
		step.Fail(err)
		return err
	}

	traceCh := make(chan []interface{})
	register := newChannelRegister(traceCh)
	go func() {
		if err := simulation.simulate(request.Method, register); err != nil {
			step.Fail(err)
		}
		register.close()
	}()

	for values := range traceCh {
		trace, err := model.NewTrace(values...)
		if err != nil {
			step.Fail(err)
			return err
		}
		err = stream.Send(&SimulateResponse{
			Trace: trace,
		})
		if err != nil {
			step.Fail(err)
			return err
		}
	}
	step.Complete()
	return nil
}
