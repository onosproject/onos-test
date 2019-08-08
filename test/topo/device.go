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
	"github.com/onosproject/onos-topo/pkg/northbound/proto"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func init() {
	test.Registry.RegisterTest("device-service", TestDeviceService, []*runner.TestSuite{TopoTests})
}

func TestDeviceService(t *testing.T) {
	conn, err := env.GetTopoConn()
	assert.NoError(t, err)
	defer conn.Close()
	client := proto.NewDeviceServiceClient(conn)

	list, err := client.List(context.Background(), &proto.ListRequest{})
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

	events := make(chan *proto.ListResponse)
	go func() {
		list, err := client.List(context.Background(), &proto.ListRequest{
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

	addResponse, err := client.Add(context.Background(), &proto.AddDeviceRequest{
		Device: &proto.Device{
			Id:      "foo",
			Address: "device-foo:5000",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "foo", addResponse.Metadata.Id)
	assert.NotEqual(t, uint64(0), addResponse.Metadata.Version)

	getResponse, err := client.Get(context.Background(), &proto.GetDeviceRequest{
		DeviceId: "foo",
	})
	assert.NoError(t, err)

	device := getResponse.Device
	assert.Equal(t, "foo", device.Id)
	assert.Equal(t, "foo", device.Metadata.Id)
	assert.Equal(t, addResponse.Metadata.Version, device.Metadata.Version)

	eventResponse := <-events
	assert.Equal(t, proto.ListResponse_ADDED, eventResponse.Type)
	assert.Equal(t, "foo", eventResponse.Device.Id)
	assert.Equal(t, "foo", eventResponse.Device.Metadata.Id)
	assert.Equal(t, addResponse.Metadata.Version, eventResponse.Device.Metadata.Version)

	list, err = client.List(context.Background(), &proto.ListRequest{})
	assert.NoError(t, err)
	for {
		response, err := list.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		assert.Equal(t, proto.ListResponse_NONE, response.Type)
		assert.Equal(t, "foo", response.Device.Id)
		assert.Equal(t, "foo", response.Device.Metadata.Id)
		assert.Equal(t, addResponse.Metadata.Version, response.Device.Metadata.Version)
		count++
	}
	assert.Equal(t, 1, count)

	removeResponse, err := client.Remove(context.Background(), &proto.RemoveDeviceRequest{
		Device: device,
	})
	assert.NoError(t, err)
	assert.NotNil(t, removeResponse)

	eventResponse = <-events
	assert.Equal(t, proto.ListResponse_REMOVED, eventResponse.Type)
	assert.Equal(t, "foo", eventResponse.Device.Id)
	assert.Equal(t, "foo", eventResponse.Device.Metadata.Id)
	assert.Equal(t, addResponse.Metadata.Version, eventResponse.Device.Metadata.Version)
}
