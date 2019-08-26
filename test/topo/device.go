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

package topo

import (
	"context"
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"github.com/onosproject/onos-test/test/env"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func init() {
	test.Registry.RegisterTest("device-service", TestDeviceService, []*runner.TestSuite{TopoTests})
}

// TestDeviceService :
func TestDeviceService(t *testing.T) {
	conn, err := env.GetTopoConn()
	assert.NoError(t, err)
	defer conn.Close()
	client := device.NewDeviceServiceClient(conn)

	list, err := client.List(context.Background(), &device.ListRequest{})
	assert.NoError(t, err)

	count := 0
	for {
		_, err := list.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		count++
	}

	assert.Equal(t, 0, count)

	events := make(chan *device.ListResponse)
	go func() {
		list, err := client.List(context.Background(), &device.ListRequest{
			Subscribe: true,
		})
		assert.NoError(t, err)

		for {
			response, err := list.Recv()
			if err == io.EOF {
				break
			}
			assert.NoError(t, err)
			events <- response
		}
	}()

	addResponse, err := client.Add(context.Background(), &device.AddRequest{
		Device: &device.Device{
			ID:      "foo",
			Address: "device-foo:5000",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, device.ID("foo"), addResponse.Device.ID)
	assert.NotEqual(t, device.Revision(0), addResponse.Device.Revision)

	getResponse, err := client.Get(context.Background(), &device.GetRequest{
		ID: "foo",
	})
	assert.NoError(t, err)

	assert.Equal(t, device.ID("foo"), getResponse.Device.ID)
	assert.Equal(t, addResponse.Device.Revision, getResponse.Device.Revision)

	eventResponse := <-events
	assert.Equal(t, device.ListResponse_ADDED, eventResponse.Type)
	assert.Equal(t, device.ID("foo"), eventResponse.Device.ID)
	assert.Equal(t, addResponse.Device.Revision, eventResponse.Device.Revision)

	list, err = client.List(context.Background(), &device.ListRequest{})
	assert.NoError(t, err)
	for {
		response, err := list.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		assert.Equal(t, device.ListResponse_NONE, response.Type)
		assert.Equal(t, device.ID("foo"), response.Device.ID)
		assert.Equal(t, addResponse.Device.Revision, response.Device.Revision)
		count++
	}
	assert.Equal(t, 1, count)

	removeResponse, err := client.Remove(context.Background(), &device.RemoveRequest{
		Device: getResponse.Device,
	})
	assert.NoError(t, err)
	assert.NotNil(t, removeResponse)

	eventResponse = <-events
	assert.Equal(t, device.ListResponse_REMOVED, eventResponse.Type)
	assert.Equal(t, device.ID("foo"), eventResponse.Device.ID)
	assert.Equal(t, addResponse.Device.Revision, eventResponse.Device.Revision)
}
