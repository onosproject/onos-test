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
	"math/rand"
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
func newSimulation(name string, suite SimulatingSuite, args map[string]string) *Simulation {
	return &Simulation{
		Name:   name,
		suite:  suite,
		args:   args,
		stopCh: make(chan error),
	}
}

// Simulation is a simulator runner
type Simulation struct {
	// Name is the name of the simulator
	Name     string
	suite    SimulatingSuite
	args     map[string]string
	register Register
	stopCh   chan error
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

// Record records an event in the register
func (s *Simulation) Record(entry interface{}) {
	s.register.Record(entry)
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

// start starts running the simulation
func (s *Simulation) start(rate time.Duration, jitter float64, register Register) {
	s.register = register
	go s.run(rate, jitter)
}

// run runs the simulation
func (s *Simulation) run(rate time.Duration, jitter float64) {
	for {
		select {
		case <-waitJitter(rate, jitter):
			s.simulateRandom()
		case <-s.stopCh:
			s.register.close()
			return
		}
	}
}

// simulateRandom calls a random simulator method
func (s *Simulation) simulateRandom() {
	method := s.chooseRandom()
	method.Func.Call([]reflect.Value{reflect.ValueOf(s.suite), reflect.ValueOf(s)})
}

// chooseRandom chooses a random simulator method
func (s *Simulation) chooseRandom() reflect.Method {
	simulators := getSimulators(s.suite)
	simulator := simulators[rand.Intn(len(simulators))]
	methods := reflect.TypeOf(s.suite)
	method, ok := methods.MethodByName(simulator)
	if !ok {
		panic(fmt.Errorf("unknown simulator method %s", simulator))
	}
	return method
}

// stop stops running the simulation
func (s *Simulation) stop() {
	close(s.stopCh)
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

// getSimulators returns a list of simulators in the given suite
func getSimulators(suite SimulatingSuite) []string {
	methodFinder := reflect.TypeOf(suite)
	simulators := []string{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := simulatorFilter(method.Name)
		if ok {
			simulators = append(simulators, method.Name)
		} else if err != nil {
			panic(err)
		}
	}
	return simulators
}

// simulatorFilter filters simulation method names
func simulatorFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Simulate", name); !ok {
		return false, nil
	}
	return true, nil
}
