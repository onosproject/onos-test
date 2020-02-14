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
	"encoding/json"
	"fmt"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/registry"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"google.golang.org/grpc"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// newWorker returns a new simulation worker
func newWorker(config *Config) (*Worker, error) {
	kubeAPI, err := kube.GetAPI(config.ID)
	if err != nil {
		return nil, err
	}
	return &Worker{
		client:      kubeAPI.Client(),
		simulations: make(map[string]*Simulation),
	}, nil
}

// Worker runs a simulation job
type Worker struct {
	client      client.Client
	simulations map[string]*Simulation
}

// Run runs a simulation
func (w *Worker) Run() error {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	RegisterSimulatorServiceServer(server, w)
	return server.Serve(lis)
}

func (w *Worker) getSimulation(name string, args map[string]string) (*Simulation, error) {
	if simulation, ok := w.simulations[name]; ok {
		return simulation, nil
	}
	suite := registry.GetSimulationSuite(name)
	if suite != nil {
		simulation := newSimulation(name, suite, args)
		w.simulations[name] = simulation
		return simulation, nil
	}
	return nil, fmt.Errorf("unknown simulation %s", name)
}

// SetupSimulation sets up a simulation suite
func (w *Worker) SetupSimulation(ctx context.Context, request *SimulationLifecycleRequest) (*SimulationLifecycleResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulationWorker()), "SetupSimulation %s", request.Simulation)
	step.Start()

	simulation, err := w.getSimulation(request.Simulation, request.Args)
	if err != nil {
		step.Fail(err)
		return nil, err
	}
	simulation.setup()
	step.Complete()
	return &SimulationLifecycleResponse{}, nil
}

// TearDownSimulation tears down a simulation suite
func (w *Worker) TearDownSimulation(ctx context.Context, request *SimulationLifecycleRequest) (*SimulationLifecycleResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulationWorker()), "TearDownSimulation %s", request.Simulation)
	step.Start()

	simulation, err := w.getSimulation(request.Simulation, request.Args)
	if err != nil {
		step.Fail(err)
		return nil, err
	}
	simulation.teardown()
	step.Complete()
	return &SimulationLifecycleResponse{}, nil
}

// StartSimulation starts a simulation
func (w *Worker) StartSimulation(request *SimulationRequest, stream SimulatorService_StartSimulationServer) error {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulationWorker()), "StartSimulation %s", request.Simulation)
	step.Start()

	simulation, err := w.getSimulation(request.Simulation, request.Args)
	if err != nil {
		step.Fail(err)
		return err
	}

	// Create a channel to read records from the register and start the simulation in the background
	traceCh := make(chan interface{})
	go simulation.start(request.Rate, request.Jitter, newChannelRegister(traceCh))
	step.Complete()

	for record := range traceCh {
		bytes, err := json.Marshal(record)
		if err != nil {
			err = stream.Send(&SimulationResponse{
				Error: err.Error(),
			})
		} else {
			err = stream.Send(&SimulationResponse{
				Result: bytes,
			})
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// StopSimulation stops a simulation
func (w *Worker) StopSimulation(ctx context.Context, request *SimulationRequest) (*SimulationResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulationWorker()), "StopSimulation %s", request.Simulation)
	step.Start()

	simulation, ok := w.simulations[request.Simulation]
	if !ok {
		return nil, fmt.Errorf("unknown simulation %s", request.Simulation)
	}
	simulation.stop()
	step.Complete()
	return &SimulationResponse{}, nil
}
