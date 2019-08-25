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
	//clusterID string
	status *console.StatusWriter
}

// Setup sets up a local cluster with a given configuration
func (c *ClusterController) Setup() console.ErrorStatus {
	c.status.Start("Setting up onos-config")

	return c.status.Succeed()
}

func (c *ClusterController) Teardown() console.ErrorStatus {
	c.status.Start("Deleting local cluster")

	return c.status.Succeed()
}
