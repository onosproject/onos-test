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

package api

import (
	"os"
	"testing"

	"github.com/onosproject/onos-test/pkg/onit/setup"
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"gotest.tools/assert"
)

func init() {
	test.Registry.RegisterTest("add-network", AddNetwork, []*runner.TestSuite{})
}

// AddNetwork test adding a stratum network to the cluster
func AddNetwork(t *testing.T) {
	clusterID := os.Getenv("TEST_NAMESPACE")
	testSetup := setup.New().
		SetClusterID(clusterID).
		SetNetworkName("stratum-1").
		Build()
	testSetup.AddNetwork()
	networks, _ := testSetup.GetNetworks()
	assert.Equal(t, len(networks), 1)
	testSetup.RemoveNetwork()
}
