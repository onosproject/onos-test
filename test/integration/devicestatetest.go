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

package integration

import (
	"context"
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"testing"
	"time"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/onosproject/onos-test/test/env"
	"github.com/stretchr/testify/assert"
)


// TestDeviceState tests that a device is connected and available.
func TestDeviceState(t *testing.T) {
	// Get the first configured device from the environment.
	envDevice := env.GetDevices()[0]
	conn, errConn := env.GetTopoConn()
	assert.NoError(t, errConn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client := device.NewDeviceServiceClient(conn)
	response, errGet := client.Get(ctx, &device.GetRequest{
		ID: device.ID(envDevice),
	})
	assert.NoError(t, errGet)
	responseDevice := response.Device
	assert.Equal(t, responseDevice.ID, device.ID(envDevice), "Wrong Device")
	assert.Equal(t, responseDevice.Protocols[0].Protocol, device.Protocol_GNMI)
	assert.Equal(t, responseDevice.Protocols[0].ConnectivityState, device.ConnectivityState_REACHABLE)
	assert.Equal(t, responseDevice.Protocols[0].ChannelState, device.ChannelState_CONNECTED)
	assert.Equal(t, responseDevice.Protocols[0].ServiceState, device.ServiceState_AVAILABLE)
}

func init() {
	test.Registry.RegisterTest("device-state", TestDeviceState, []*runner.TestSuite{AllTests, IntegrationTests})
}
