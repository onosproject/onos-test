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

package local

import "github.com/onosproject/onos-test/pkg/onit/console"

// ClusterController manages a single local cluster
type ClusterController struct {
	clusterID string
	status    *console.StatusWriter
}

// Teardown tears down a local cluster
func (c *ClusterController) Teardown() console.ErrorStatus {

	c.status.Start("Tearing down the local cluster")
	return c.status.Succeed()
}

// Setup sets up the local cluster
func (c *ClusterController) Setup() console.ErrorStatus {

	c.status.Start("Setting up local Atomix controller")
	err := c.setupAtomixController()
	if err != nil {
		return c.status.Fail(err)
	}
	c.status.Succeed()

	c.status.Start("Bootstrapping local onos-config cluster")
	err = c.setupOnosConfig()
	if err != nil {
		return c.status.Fail(err)
	}
	c.status.Succeed()

	c.status.Start("Bootstrapping local onos-topo cluster")
	err = c.setupOnosTopo()
	if err != nil {
		return c.status.Fail(err)
	}
	c.status.Succeed()

	return c.status.Succeed()
}
