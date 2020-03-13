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
	"strconv"
	"sync"
	"time"
)

// SimulatingSuite is a suite of simulators
type SimulatingSuite interface{}

// Suite is an identifier interface for simulation suites
type Suite struct{}

// ScheduleSimulator is an interface for scheduling operations for a simulation
type ScheduleSimulator interface {
	ScheduleSimulator(s *Simulator)
}

// SetupSimulation is an interface for setting up a suite of simulators
type SetupSimulation interface {
	SetupSimulation(s *Simulator) error
}

// TearDownSimulation is an interface for tearing down a suite of simulators
type TearDownSimulation interface {
	TearDownSimulation(s *Simulator) error
}

// SetupSimulator is an interface for executing code before every simulator
type SetupSimulator interface {
	SetupSimulator(s *Simulator) error
}

// TearDownSimulator is an interface for executing code after every simulator
type TearDownSimulator interface {
	TearDownSimulator(s *Simulator) error
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

// newSimulator returns a new simulation instance
func newSimulator(name string, process int, suite SimulatingSuite, args map[string]string, config *Config) *Simulator {
	return &Simulator{
		Name:    name,
		Process: process,
		suite:   suite,
		args:    args,
		config:  config,
		ops:     make(map[string]*operation),
	}
}

// Simulator is a simulator runner
type Simulator struct {
	// Name is the name of the simulation
	Name string
	// Process is the unique identifier of the simulator process
	Process int
	config  *Config
	suite   SimulatingSuite
	args    map[string]string
	ops     map[string]*operation
	mu      sync.Mutex
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

// Schedule schedules an operation
func (s *Simulator) Schedule(name string, f func(*Simulator) error, rate time.Duration, jitter float64) {
	if override, ok := s.config.Rates[name]; ok {
		rate = override
	}
	if override, ok := s.config.Jitter[name]; ok {
		jitter = override
	}
	s.ops[name] = &operation{
		name:      name,
		f:         f,
		rate:      rate,
		jitter:    jitter,
		simulator: s,
		stopCh:    make(chan error),
	}
}

// lock locks the simulation
func (s *Simulator) lock() {
	s.mu.Lock()
}

// unlock unlocks the simulator
func (s *Simulator) unlock() {
	s.mu.Unlock()
}

// setupSimulation sets up the simulation
func (s *Simulator) setupSimulation() error {
	if setupSuite, ok := s.suite.(SetupSimulation); ok {
		return setupSuite.SetupSimulation(s)
	}
	return nil
}

// teardownSimulation tears down the simulation
func (s *Simulator) teardownSimulation() error {
	if tearDownSuite, ok := s.suite.(TearDownSimulation); ok {
		return tearDownSuite.TearDownSimulation(s)
	}
	return nil
}

// setupSimulator sets up the simulator
func (s *Simulator) setupSimulator() error {
	if setupSuite, ok := s.suite.(SetupSimulator); ok {
		if err := setupSuite.SetupSimulator(s); err != nil {
			return err
		}
	}
	if setupSuite, ok := s.suite.(ScheduleSimulator); ok {
		setupSuite.ScheduleSimulator(s)
	}
	return nil
}

// teardownSimulator tears down the simulator
func (s *Simulator) teardownSimulator() error {
	if tearDownSuite, ok := s.suite.(TearDownSimulator); ok {
		return tearDownSuite.TearDownSimulator(s)
	}
	return nil
}

// start starts the simulator
func (s *Simulator) start() {
	for _, op := range s.ops {
		go op.start()
	}
}

// stop stops the simulator
func (s *Simulator) stop() {
	for _, op := range s.ops {
		op.stop()
	}
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

// operation is a simulator operation
type operation struct {
	name      string
	f         func(*Simulator) error
	rate      time.Duration
	jitter    float64
	simulator *Simulator
	stopCh    chan error
}

// start starts the operation simulator
func (o *operation) start() {
	for {
		select {
		case <-waitJitter(o.rate, o.jitter):
			o.simulator.lock()
			o.run()
			o.simulator.unlock()
		case <-o.stopCh:
			return
		}
	}
}

// run runs the operation
func (o *operation) run() {
	step := logging.NewStep(fmt.Sprintf("%s/%d", o.simulator.Name, getSimulatorID()), "Run %s", o.name)
	step.Start()
	if err := o.f(o.simulator); err != nil {
		step.Fail(err)
	} else {
		step.Complete()
	}
}

// stop stops the operation simulator
func (o *operation) stop() {
	close(o.stopCh)
}

// newSimulatorServer returns a new simulator server
func newSimulatorServer(config *Config) (*simulatorServer, error) {
	return &simulatorServer{
		config:      config,
		simulations: make(map[string]*Simulator),
	}, nil
}

// simulatorServer listens for simulator requests
type simulatorServer struct {
	config      *Config
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
		simulation := newSimulator(name, getSimulatorID(), suite, args, s.config)
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
	if err := simulation.setupSimulation(); err != nil {
		step.Fail(err)
		return nil, err
	}
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
	if err := simulation.teardownSimulation(); err != nil {
		step.Fail(err)
		return nil, err
	}
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
	if err := simulation.setupSimulator(); err != nil {
		step.Fail(err)
		return nil, err
	}
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
	if err := simulation.teardownSimulator(); err != nil {
		step.Fail(err)
		return nil, err
	}
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

	go simulation.start()
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
