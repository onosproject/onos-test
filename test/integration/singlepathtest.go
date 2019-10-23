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
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"strconv"
	"testing"

	"github.com/onosproject/onos-test/test/env"
	"github.com/stretchr/testify/assert"
)

const (
	tzValue = "Europe/Dublin"
	tzPath  = "/system/clock/config/timezone-name"
)

func makeDevicePath(device string, path string) []DevicePath {
	devicePath := make([]DevicePath, 1)
	devicePath[0].deviceName = device
	devicePath[0].path = path
	return devicePath
}

// TestSinglePath tests query/set/delete of a single GNMI path to a single device
func TestSinglePath(t *testing.T) {
	// Get the first configured device from the environment.
	device := env.GetDevices()[0]

	// Make a GNMI client to use for requests
	c, err := env.NewGnmiClient(MakeContext(), "")
	assert.NoError(t, err)
	assert.True(t, c != nil, "Fetching client returned nil")

	// Set a value using gNMI client
	setPath := makeDevicePath(device, tzPath)
	setPath[0].pathDataValue = tzValue
	setPath[0].pathDataType = StringVal
	_, extensions, errorSet := GNMISet(MakeContext(), c, setPath, noPaths)
	assert.NoError(t, errorSet)
	assert.Equal(t, 1, len(extensions))
	extension := extensions[0].GetRegisteredExt()
	assert.Equal(t, extension.Id.String(), strconv.Itoa(100))

	// Check that the value was set correctly
	valueAfter, extensions, errorAfter := GNMIGet(MakeContext(), c, makeDevicePath(device, tzPath))
	assert.NoError(t, errorAfter)
	assert.Equal(t, 0, len(extensions))
	assert.NotEqual(t, "", valueAfter, "Query after set returned an error: %s\n", errorAfter)
	assert.Equal(t, tzValue, valueAfter[0].pathDataValue, "Query after set returned the wrong value: %s\n", valueAfter)

	// Remove the path we added
	_, extensions, errorDelete := GNMISet(MakeContext(), c, noPaths, makeDevicePath(device, tzPath))
	assert.NoError(t, errorDelete)
	assert.Equal(t, 1, len(extensions))
	extension = extensions[0].GetRegisteredExt()
	assert.Equal(t, extension.Id.String(), strconv.Itoa(100))

	//  Make sure it got removed
	valueAfterDelete, extensions, errorAfterDelete := GNMIGet(MakeContext(), c, makeDevicePath(device, tzPath))
	assert.NoError(t, errorAfterDelete)
	assert.Equal(t, 0, len(extensions))
	assert.Equal(t, valueAfterDelete[0].pathDataValue, "",
		"incorrect value found for path /system/clock/config/timezone-name after delete")
}

func init() {
	test.Registry.RegisterTest("single-path", TestSinglePath, []*runner.TestSuite{AllTests, IntegrationTests})
}
