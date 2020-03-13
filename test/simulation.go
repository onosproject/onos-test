package test

import (
	"github.com/onosproject/onos-test/pkg/helm"
	"github.com/onosproject/onos-test/pkg/simulation"
	"time"
)

type ChartSimulationSuite struct {
	simulation.Suite
}

// SetupSimulation :: simulation
func (s *ChartSimulationSuite) SetupSimulation(sim *simulation.Simulator) error {
	return helm.Helm().
		Chart("/etc/charts/atomix-controller").
		Release("atomix-controller").
		Set("namespace", helm.Namespace()).
		Install(true)
}

// ScheduleSimulator :: simulation
func (s *ChartSimulationSuite) ScheduleSimulator(sim *simulation.Simulator) {
	sim.Schedule("foo", s.SimulateFoo, 1*time.Second, 1)
	sim.Schedule("bar", s.SimulateBar, 5*time.Second, 1)
	sim.Schedule("baz", s.SimulateBaz, 30*time.Second, 1)
}

func (s *ChartSimulationSuite) SimulateFoo(sim *simulation.Simulator) error {
	println(sim.Arg("foo").String("<none>"))
	return nil
}

func (s *ChartSimulationSuite) SimulateBar(sim *simulation.Simulator) error {
	println(sim.Arg("bar").String("<none>"))
	return nil
}

func (s *ChartSimulationSuite) SimulateBaz(sim *simulation.Simulator) error {
	println(sim.Arg("baz").String("<none>"))
	return nil
}
