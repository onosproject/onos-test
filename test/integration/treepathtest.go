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
	newRootName            = "new-root"
	newRootPath            = "/interfaces/interface[name=" + newRootName + "]"
	newRootConfigNamePath  = newRootPath + "/config/name"
	newRootEnabledPath     = newRootPath + "/config/enabled"
	newRootDescriptionPath = newRootPath + "/config/description"
	newDescription         = "description"
)

// TestTreePath tests create/set/delete of a tree of GNMI paths to a single device
func TestTreePath(t *testing.T) {
	// Get the first configured device from the environment.
	device := env.GetDevices()[0]

	// Make a GNMI client to use for requests
	c, err := env.NewGnmiClient(MakeContext(), "")
	assert.NoError(t, err)
	assert.True(t, c != nil, "Fetching client returned nil")

	// Set name of new root using gNMI client
	setNamePath := []DevicePath{
		{deviceName: device, path: newRootConfigNamePath, pathDataValue: newRootName, pathDataType: StringVal},
	}
	_, errorSet := GNMISet(MakeContext(), c, setNamePath)
	assert.NoError(t, errorSet)

	// Set values using gNMI client
	setPath := []DevicePath{
		{deviceName: device, path: newRootDescriptionPath, pathDataValue: newDescription, pathDataType: StringVal},
		{deviceName: device, path: newRootEnabledPath, pathDataValue: "false", pathDataType: BoolVal},
	}
	_, errorSet = GNMISet(MakeContext(), c, setPath)
	assert.NoError(t, errorSet)

	// Check that the name value was set correctly
	valueAfter, errorAfter := GNMIGet(MakeContext(), c, setNamePath)
	assert.NoError(t, errorAfter)
	assert.NotEqual(t, "", valueAfter, "Query name after set returned an error: %s\n", errorAfter)
	assert.Equal(t, newRootName, valueAfter[0].pathDataValue, "Query name after set returned the wrong value: %s\n", valueAfter)

	// Check that the enabled value was set correctly
	valueAfter, errorAfter = GNMIGet(MakeContext(), c, makeDevicePath(device, newRootEnabledPath))
	assert.NoError(t, errorAfter)
	assert.NotEqual(t, "", valueAfter, "Query enabled after set returned an error: %s\n", errorAfter)
	assert.Equal(t, "false", valueAfter[0].pathDataValue, "Query enabled after set returned the wrong value: %s\n", valueAfter)

	// Remove the root path we added
	errorDelete := GNMIDelete(MakeContext(), c, makeDevicePath(device, newRootPath))
	assert.NoError(t, errorDelete)

	//  Make sure child got removed
	valueAfterDelete, errorAfterDelete := GNMIGet(MakeContext(), c, makeDevicePath(device, newRootConfigNamePath))
	assert.NoError(t, errorAfterDelete)
	assert.Equal(t, valueAfterDelete[0].pathDataValue, "",
		"New child was not removed")

	//  Make sure new root got removed
	valueAfterRootDelete, errorAfterRootDelete := GNMIGet(MakeContext(), c, makeDevicePath(device, newRootPath))
	assert.NoError(t, errorAfterRootDelete)
	assert.Equal(t, valueAfterRootDelete[0].pathDataValue, "",
		"New root was not removed")
}

func init() {
	test.Registry.RegisterTest("tree-path", TestTreePath, []*runner.TestSuite{AllTests, IntegrationTests})
}
