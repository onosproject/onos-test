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
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// SimulatingSuite is a suite of simulators
type SimulatingSuite interface{}

// Suite is an identifier interface for simulation suites
type Suite struct{}

// SetupSimulation is an interface for setting up a suite of simulators
type SetupSimulation interface {
	SetupSimulation(s *Simulation)
}

// TearDownSimulation is an interface for tearing down a suite of simulators
type TearDownSimulation interface {
	TearDownSimulation(s *Simulation)
}

// SetupSimulator is an interface for executing code before every simulator
type SetupSimulator interface {
	SetupSimulator(s *Simulation)
}

// TearDownSimulator is an interface for executing code after every simulator
type TearDownSimulator interface {
	TearDownSimulator(s *Simulation)
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
func newSimulation(name string, process int, suite SimulatingSuite, args map[string]string) *Simulation {
	return &Simulation{
		Name:    name,
		Process: process,
		suite:   suite,
		args:    args,
	}
}

// Simulation is a simulator runner
type Simulation struct {
	// Name is the name of the simulation
	Name string
	// Process is the unique identifier of the simulator process
	Process  int
	suite    SimulatingSuite
	args     map[string]string
	register Register
}

// withRegister returns an instance of the simulator with the given register
func (s *Simulation) withRegister(register Register) *Simulation {
	return &Simulation{
		Name:     s.Name,
		Process:  s.Process,
		suite:    s.suite,
		args:     s.args,
		register: register,
	}
}

// Arg gets a simulator argument
func (s *Simulation) Arg(name string) *Arg {
	if value, ok := s.args[name]; ok {
		return &Arg{
			value: value,
		}
	}
	return &Arg{}
}

// Trace records an trace in the register
func (s *Simulation) Trace(values ...interface{}) {
	s.register.Trace(values...)
}

// setup sets up the simulation
func (s *Simulation) setup() {
	s.setupSimulation()
}

// teardown tears down the simulation
func (s *Simulation) teardown() {
	s.teardownSimulation()
}

// setupSimulation sets up the simulation
func (s *Simulation) setupSimulation() {
	if setupSuite, ok := s.suite.(SetupSimulation); ok {
		setupSuite.SetupSimulation(s)
	}
}

// teardownSimulation tears down the simulation
func (s *Simulation) teardownSimulation() {
	if tearDownSuite, ok := s.suite.(TearDownSimulation); ok {
		tearDownSuite.TearDownSimulation(s)
	}
}

// setupSimulator sets up the simulator
func (s *Simulation) setupSimulator() {
	if setupSuite, ok := s.suite.(SetupSimulator); ok {
		setupSuite.SetupSimulator(s)
	}
}

// teardownSimulator tears down the simulator
func (s *Simulation) teardownSimulator() {
	if tearDownSuite, ok := s.suite.(TearDownSimulator); ok {
		tearDownSuite.TearDownSimulator(s)
	}
}

// simulate simulates the given method
func (s *Simulation) simulate(name string, register Register) error {
	methods := reflect.TypeOf(s.suite)
	method, ok := methods.MethodByName(name)
	if !ok {
		return fmt.Errorf("unknown simulator method %s", name)
	}
	method.Func.Call([]reflect.Value{reflect.ValueOf(s.suite), reflect.ValueOf(s.withRegister(register))})
	return nil
}

// simulatorFilter filters simulation method names
func simulatorFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Simulate", name); !ok {
		return false, nil
	}
	return true, nil
}
