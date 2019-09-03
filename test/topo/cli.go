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
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"github.com/onosproject/onos-test/test/env"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func init() {
	test.Registry.RegisterTest("device-cli", TestTopoDeviceCLI, []*runner.TestSuite{TopoTests})
}

const (
	devicesHeader = "ID ADDRESS VERSION"
	addedTest     = "Added device test"
	removedTest   = "Removed device test"
	devicesTest   = "test test:1234 1.0.0"
)

// TestTopoDeviceCLI tests the topo service's device CLI commands
func TestTopoDeviceCLI(t *testing.T) {
	var output []string
	var code int

	output, code = env.ExecuteCLI("onos topo get devices")
	assert.Equal(t, 0, code)
	assert.Len(t, output, 1)
	assert.Equal(t, devicesHeader, stripSpaces(output[0]))

	output, code = env.ExecuteCLI("onos topo add device test --type Devicesim --address test:1234 --version 1.0.0")
	assert.Equal(t, 0, code)
	assert.Len(t, output, 1)
	assert.Equal(t, addedTest, output[0])

	output, code = env.ExecuteCLI("onos topo get devices")
	assert.Equal(t, 0, code)
	assert.Len(t, output, 2)
	assert.Equal(t, devicesHeader, stripSpaces(output[0]))
	assert.Equal(t, devicesTest, stripSpaces(output[1]))

	output, code = env.ExecuteCLI("onos topo remove device test")
	assert.Equal(t, 0, code)
	assert.Len(t, output, 1)
	assert.Equal(t, removedTest, output[0])

	output, code = env.ExecuteCLI("onos topo get devices")
	assert.Equal(t, 0, code)
	assert.Len(t, output, 1)
	assert.Equal(t, devicesHeader, stripSpaces(output[0]))
}

func stripSpaces(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, " ")
}
