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

import (
	"github.com/onosproject/onos-test/pkg/onit/console"
)

// NewController creates a new local controller
func NewController() (*Controller, error) {

	return &Controller{
		console.NewStatusWriter(),
	}, nil

}

// Controller is a local controller that manages local clusters for onit
type Controller struct {
	status *console.StatusWriter
}

// GetClusters returns a list of local clusters
func (c *Controller) GetClusters() (map[string]*ClusterConfig, error) {

	return nil, nil
}

// NewClusterController creates a new instance of local ClusterController
func (c *Controller) NewClusterController(clusterID string, config *ClusterConfig) *ClusterController {
	return &ClusterController{
		clusterID: clusterID,
		status:    c.status,
	}
}

// NewCluster creates a new local cluster controller
func (c *Controller) NewCluster(clusterID string, config *ClusterConfig) (*ClusterController, console.ErrorStatus) {

	return c.NewClusterController(clusterID, config), c.status.Succeed()
}

// GetCluster returns a local cluster controller
func (c *Controller) GetCluster(clusterID string) (*ClusterController, error) {

	return nil, nil
}
