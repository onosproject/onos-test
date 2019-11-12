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

package config

import (
	"context"
	"github.com/onosproject/onos-test/pkg/onit"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testSuite struct {
	onit.TestSuite
}

// addSimulator adds a device to the network
func (s *testSuite) addSimulator(t *testing.T) env.SimulatorEnv {
	simulator := env.Simulators().New().AddOrDie()

	conn, err := env.Topo().Connect()
	assert.NoError(t, err)

	client := device.NewDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	_, err = client.Add(ctx, &device.AddRequest{
		Device: &device.Device{
			ID:      device.ID(simulator.Name()),
			Address: simulator.Address(),
			Type:    "Simulator",
			Version: "1.0.0",
		},
	})
	cancel()
	assert.NoError(t, err)
	return simulator
}

// SmokeTestSuite is the primary onos-config test suite
type SmokeTestSuite struct {
	testSuite
}

// SetupTestSuite sets up the onos-config test suite
func (s *SmokeTestSuite) SetupTestSuite() {
	setup.Topo().Nodes(2)
	setup.Config().Nodes(2)
	setup.SetupOrDie()
}

// CLITestSuite is the onos-config CLI test suite
type CLITestSuite struct {
	testSuite
}

// SetupTestSuite sets up the onos-config CLI test suite
func (s *CLITestSuite) SetupTestSuite() {
	setup.Topo().Nodes(2)
	setup.Config().Nodes(2)
	setup.SetupOrDie()
}

// HATestSuite is the onos-config HA test suite
type HATestSuite struct {
	testSuite
}

// SetupTestSuite sets up the onos-config CLI test suite
func (s *HATestSuite) SetupTestSuite() {
	setup.Topo().Nodes(2)
	setup.Config().Nodes(2)
	setup.SetupOrDie()
}
