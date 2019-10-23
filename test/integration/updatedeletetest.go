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
	"testing"

	"github.com/onosproject/onos-test/test/env"
	"github.com/stretchr/testify/assert"
)

const (
	udtestRootPath         = "/interfaces/interface[name=test]"
	udtestNamePath         = udtestRootPath + "/config/name"
	udtestEnabledPath      = udtestRootPath + "/config/enabled"
	udtestDescriptionPath  = udtestRootPath + "/config/description"
	udtestNameValue        = "test"
	udtestDescriptionValue = "description"
)

// TestUpdateDelete tests update and delete paths in a single GNMI request
func TestUpdateDelete(t *testing.T) {
	// Get the first configured device from the environment.
	device := env.GetDevices()[0]

	// Make a GNMI client to use for requests
	c, err := env.NewGnmiClient(MakeContext(), "")
	assert.NoError(t, err)
	assert.True(t, c != nil, "Fetching client returned nil")

	noPaths := make([]DevicePath, 0)

	// Create interface tree using gNMI client
	setNamePath := []DevicePath{
		{deviceName: device, path: udtestNamePath, pathDataValue: udtestNameValue, pathDataType: StringVal},
	}
	_, _, errorSet := GNMIUpdateAndDelete(MakeContext(), c, setNamePath, noPaths)
	assert.NoError(t, errorSet)

	// Set initial values for Enabled and Description using gNMI client
	setInitialValuesPath := []DevicePath{
		{deviceName: device, path: udtestEnabledPath, pathDataValue: "true", pathDataType: BoolVal},
		{deviceName: device, path: udtestDescriptionPath, pathDataValue: udtestDescriptionValue, pathDataType: StringVal},
	}
	_, _, errorSet = GNMIUpdateAndDelete(MakeContext(), c, setInitialValuesPath, noPaths)
	assert.NoError(t, errorSet)

	// Update Enabled, delete Description using gNMI client
	updateEnabledPath := []DevicePath{
		{deviceName: device, path: udtestEnabledPath, pathDataValue: "false", pathDataType: BoolVal},
	}
	deleteDescriptionPath := []DevicePath{
		{deviceName: device, path: udtestDescriptionPath},
	}
	_, _, errorSet = GNMIUpdateAndDelete(MakeContext(), c, updateEnabledPath, deleteDescriptionPath)
	assert.NoError(t, errorSet)

	// Check that the Enabled value is set correctly
	valueAfter, extensions, errorAfter := GNMIGet(MakeContext(), c, updateEnabledPath)
	assert.NoError(t, errorAfter)
	assert.Equal(t, 0, len(extensions))
	assert.NotEqual(t, "", valueAfter, "Query name after set returned an error: %s\n", errorAfter)
	assert.Equal(t, "false", valueAfter[0].pathDataValue, "Query name after set returned the wrong value: %s\n", valueAfter)

	//  Make sure Description got removed
	valueAfterDelete, extensions, errorAfterDelete := GNMIGet(MakeContext(), c, makeDevicePath(device, udtestDescriptionPath))
	assert.NoError(t, errorAfterDelete)
	assert.Equal(t, valueAfterDelete[0].pathDataValue, "", "New child was not removed")
	assert.Equal(t, 0, len(extensions))
}

func init() {
	test.Registry.RegisterTest("update-delete", TestUpdateDelete, []*runner.TestSuite{AllTests, IntegrationTests})
}
