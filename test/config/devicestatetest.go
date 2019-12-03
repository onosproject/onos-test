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
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/onosproject/onos-topo/api/device"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

// TestDeviceState tests that a device is connected and available.
func (s *SmokeTestSuite) TestDeviceState(t *testing.T) {
	simulator := env.NewSimulator().AddOrDie()
	waitForConnected(device.ID(simulator.Name()), t)
}

// waitForConnected waits for the given device to connect
func waitForConnected(id device.ID, t *testing.T) {
	conn, err := env.Topo().Connect()
	assert.NoError(t, err)
	client := device.NewDeviceServiceClient(conn)

	// Set a timer within which the device must reach the connected/available state
	timer := time.NewTimer(5 * time.Second)
	go func() {
		_, ok := <-timer.C
		if !ok {
			t.Fail()
		}
	}()

	// Open a stream to listen for events from the device service
	stream, err := client.List(context.Background(), &device.ListRequest{
		Subscribe: true,
	})
	assert.NoError(t, err)

	// Wait for a device event indicating the device is connected/available
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Error(err)
		}

		responseDevice := response.Device
		assert.Equal(t, responseDevice.ID, id, "Wrong Device")
		if len(responseDevice.Protocols) > 0 &&
			responseDevice.Protocols[0].Protocol == device.Protocol_GNMI &&
			responseDevice.Protocols[0].ConnectivityState == device.ConnectivityState_REACHABLE &&
			responseDevice.Protocols[0].ChannelState == device.ChannelState_CONNECTED &&
			responseDevice.Protocols[0].ServiceState == device.ServiceState_AVAILABLE {
			timer.Stop()
			break
		}
	}
}
