package test

import (
	"github.com/onosproject/onos-test/pkg/onit/setup2"
	"testing"
)

// Scenario is a test scenario
type Scenario interface {
	// SetUp sets up the scenario
	SetUp(s setup2.TestSetup) error

	// Run runs the scenario
	Run(t *testing.T, e setup2.TestEnv)

	// TearDown tears down the test scenario
	TearDown(s setup2.TestSetup) error
}
