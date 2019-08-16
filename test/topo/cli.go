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
	addedFoo = "Added device foo"
	removedFoo = "Removed device foo"
	devicesFoo = "foo foo:1234 1.0.0"
)

// TestTopoDeviceCLI tests the topo service's device CLI commands
func TestTopoDeviceCLI(t *testing.T) {
	var output []string
	var code int

	output, code = env.ExecuteCLI("/bin/bash", "-c", "onos topo get devices")
	assert.Equal(t, 0, code)
	assert.Len(t, output, 1)
	assert.Equal(t, devicesHeader, stripSpaces(output[0]))

	output, code = env.ExecuteCLI("/bin/bash", "-c", "onos topo add device foo --type Devicesim --address foo:1234 --version 1.0.0")
	assert.Equal(t, 0, code)
	assert.Len(t, output, 1)
	assert.Equal(t, addedFoo, output[0])

	output, code = env.ExecuteCLI("/bin/bash", "-c", "onos topo get devices")
	assert.Equal(t, 0, code)
	assert.Len(t, output, 2)
	assert.Equal(t, devicesHeader, stripSpaces(output[0]))
	assert.Equal(t, devicesFoo, stripSpaces(output[1]))

	output, code = env.ExecuteCLI("/bin/bash", "-c", "onos topo remove device foo")
	assert.Equal(t, 0, code)
	assert.Len(t, output, 1)
	assert.Equal(t, removedFoo, output[0])

	output, code = env.ExecuteCLI("/bin/bash", "-c", "onos topo get devices")
	assert.Equal(t, 0, code)
	assert.Len(t, output, 1)
	assert.Equal(t, devicesHeader, stripSpaces(output[0]))
}

func stripSpaces(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, " ")
}
