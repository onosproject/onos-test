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

package cluster

import (
	"github.com/onosproject/onos-test/pkg/kube"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

// New returns a new onit Env
func New(api kube.API) *Cluster {
	objectsClient := metav1.NewObjectsClient(api, func(meta metav1.Object) (bool, error) {
		return true, nil
	})
	cluster := &Cluster{
		Client: newClient(objectsClient),
		API:    api,
		charts: make(map[string]*Chart),
	}
	return cluster
}

var _ Client = &Cluster{}

// Cluster facilitates modifying subsystems in Kubernetes
type Cluster struct {
	Client
	kube.API
	charts map[string]*Chart
}

// Chart returns a chart
func (c *Cluster) Chart(name string) *Chart {
	chart, ok := c.charts[name]
	if !ok {
		chart = newChart(name, c)
		c.charts[name] = chart
	}
	return chart
}
