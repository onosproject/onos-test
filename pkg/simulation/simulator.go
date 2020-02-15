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
	"github.com/onosproject/onos-test/pkg/registry"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

// SimulatingSuite is a suite of simulators
type SimulatingSuite interface{}

// Suite is an identifier interface for simulation suites
type Suite struct{}

// SetupSimulation is an interface for setting up a suite of simulators
type SetupSimulation interface {
	SetupSimulation(s *Simulator)
}

// TearDownSimulation is an interface for tearing down a suite of simulators
type TearDownSimulation interface {
	TearDownSimulation(s *Simulator)
}

// SetupSimulator is an interface for executing code before every simulator
type SetupSimulator interface {
	SetupSimulator(s *Simulator)
}

// TearDownSimulator is an interface for executing code after every simulator
type TearDownSimulator interface {
	TearDownSimulator(s *Simulator)
}

// Arg is a simulator argument
type Arg struct {
	value string
}

// Int returns the argument as an int
func (a *Arg) Int(def int) int {
	if a.value == "" {
		return def
	}
	i, err := strconv.Atoi(a.value)
	if err != nil {
		panic(err)
	}
	return i
}

// String returns the argument as a string
func (a *Arg) String(def string) string {
	if a.value == "" {
		return def
	}
	return a.value
}

// newSimulation returns a new simulation instance
func newSimulation(name string, process int, suite SimulatingSuite, args map[string]string) *Simulator {
	return &Simulator{
		Name:    name,
		Process: process,
		suite:   suite,
		args:    args,
		stopCh:  make(chan error),
	}
}

// Simulator is a simulator runner
type Simulator struct {
	// Name is the name of the simulation
	Name string
	// Process is the unique identifier of the simulator process
	Process  int
	suite    SimulatingSuite
	args     map[string]string
	register Register
	stopCh   chan error
}

// Arg gets a simulator argument
func (s *Simulator) Arg(name string) *Arg {
	if value, ok := s.args[name]; ok {
		return &Arg{
			value: value,
		}
	}
	return &Arg{}
}

// Trace records an trace in the register
func (s *Simulator) Trace(values ...interface{}) {
	s.register.Trace(values...)
}

// setup sets up the simulation
func (s *Simulator) setup() {
	s.setupSimulation()
}

// teardown tears down the simulation
func (s *Simulator) teardown() {
	s.teardownSimulation()
}

// setupSimulation sets up the simulation
func (s *Simulator) setupSimulation() {
	if setupSuite, ok := s.suite.(SetupSimulation); ok {
		setupSuite.SetupSimulation(s)
	}
}

// teardownSimulation tears down the simulation
func (s *Simulator) teardownSimulation() {
	if tearDownSuite, ok := s.suite.(TearDownSimulation); ok {
		tearDownSuite.TearDownSimulation(s)
	}
}

// setupSimulator sets up the simulator
func (s *Simulator) setupSimulator() {
	if setupSuite, ok := s.suite.(SetupSimulator); ok {
		setupSuite.SetupSimulator(s)
	}
}

// teardownSimulator tears down the simulator
func (s *Simulator) teardownSimulator() {
	if tearDownSuite, ok := s.suite.(TearDownSimulator); ok {
		tearDownSuite.TearDownSimulator(s)
	}
}

// start starts the simulator
func (s *Simulator) start(rate time.Duration, jitter float64, register Register) {
	s.register = register
	go s.run(rate, jitter)
}

// run runs the simulator
func (s *Simulator) run(rate time.Duration, jitter float64) {
	for {
		select {
		case <-waitJitter(rate, jitter):
			s.simulate()
		case <-s.stopCh:
			s.register.close()
			return
		}
	}
}

// stop stops the simulator
func (s *Simulator) stop() {
	close(s.stopCh)
}

// simulate simulates a random simulator method
func (s *Simulator) simulate() {
	method := s.chooseMethod()
	method.Func.Call([]reflect.Value{reflect.ValueOf(s.suite), reflect.ValueOf(s)})
}

// chooseMethod chooses a random method
func (s *Simulator) chooseMethod() reflect.Method {
	suiteType := reflect.TypeOf(s.suite)
	methods := make(map[string]reflect.Method)
	names := []string{}
	for index := 0; index < suiteType.NumMethod(); index++ {
		method := suiteType.Method(index)
		ok, err := methodFilter(method.Name)
		if ok {
			methods[method.Name] = method
			names = append(names, method.Name)
		} else if err != nil {
			panic(err)
		}
	}
	return methods[names[rand.Intn(len(names))]]
}

// waitJitter returns a channel that closes after time.Duration between duration and duration + maxFactor *
// duration.
func waitJitter(duration time.Duration, maxFactor float64) <-chan time.Time {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	delay := duration + time.Duration(rand.Float64()*maxFactor*float64(duration))
	return time.After(delay)
}

// methodFilter filters simulation method names
func methodFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Simulate", name); !ok {
		return false, nil
	}
	return true, nil
}

// newSimulatorServer returns a new simulator server
func newSimulatorServer() (*simulatorServer, error) {
	return &simulatorServer{
		simulations: make(map[string]*Simulator),
	}, nil
}

// simulatorServer listens for simulator requests
type simulatorServer struct {
	simulations map[string]*Simulator
}

// Run runs a simulation
func (s *simulatorServer) Run() error {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	RegisterSimulatorServiceServer(server, s)
	return server.Serve(lis)
}

func (s *simulatorServer) getSimulation(name string, args map[string]string) (*Simulator, error) {
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
func (s *simulatorServer) SetupSimulation(ctx context.Context, request *SimulationLifecycleRequest) (*SimulationLifecycleResponse, error) {
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
func (s *simulatorServer) TearDownSimulation(ctx context.Context, request *SimulationLifecycleRequest) (*SimulationLifecycleResponse, error) {
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

func (s *simulatorServer) SetupSimulator(ctx context.Context, request *SimulationLifecycleRequest) (*SimulationLifecycleResponse, error) {
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

func (s *simulatorServer) TearDownSimulator(ctx context.Context, request *SimulationLifecycleRequest) (*SimulationLifecycleResponse, error) {
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

func (s *simulatorServer) StartSimulator(ctx context.Context, request *SimulatorRequest) (*SimulatorResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulatorID()), "StartSimulator %s", request.Simulation)
	step.Start()

	simulation, ok := s.simulations[request.Simulation]
	if !ok {
		err := fmt.Errorf("unknown simulation %s", request.Simulation)
		step.Fail(err)
		return nil, err
	}

	register, err := newBlockingRegister(request.Register)
	if err != nil {
		return nil, err
	}

	go simulation.start(request.Rate, request.Jitter, register)
	step.Complete()
	return &SimulatorResponse{}, nil
}

func (s *simulatorServer) StopSimulator(ctx context.Context, request *SimulatorRequest) (*SimulatorResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Simulation, getSimulatorID()), "StopSimulator %s", request.Simulation)
	step.Start()

	simulation, ok := s.simulations[request.Simulation]
	if !ok {
		err := fmt.Errorf("unknown simulation %s", request.Simulation)
		step.Fail(err)
		return nil, err
	}

	simulation.stop()
	step.Complete()
	return &SimulatorResponse{}, nil
}
