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

package interfaces

import (
	"github.com/onosproject/onos-test/pkg/onit/console"
	"github.com/onosproject/onos-test/pkg/onit/k8s"
)

// ClusterController interface
type ClusterController interface {
	Setup() console.ErrorStatus
	Teardown() console.ErrorStatus
}

// Controller interface
type Controller interface {
	NewClusterController(clusterID string, config *k8s.ClusterConfig) *k8s.ClusterController
	NewCluster(string, *k8s.ClusterConfig) (*k8s.ClusterController, console.ErrorStatus)
	GetCluster(string) (*k8s.ClusterController, error)
	GetClusters() (map[string]*k8s.ClusterConfig, error)
}
