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
	"helm.sh/helm/v3/pkg/cli"
)

// New returns a new onit Env
func New(api kube.API) *Cluster {
	objectsClient := metav1.NewObjectsClient(api, func(meta metav1.Object) (bool, error) {
		return true, nil
	})
	client := newClient(objectsClient)
	cluster := &Cluster{
		Client: client,
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

// Charts returns a list of charts in the cluster
func (c *Cluster) Charts() []*Chart {
	charts := make([]*Chart, 0, len(c.charts))
	for _, chart := range c.charts {
		charts = append(charts, chart)
	}
	return charts
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

// Releases returns a list of releases
func (c *Cluster) Releases() []*Release {
	releases := make([]*Release, 0)
	for _, chart := range c.charts {
		for _, release := range chart.Releases() {
			releases = append(releases, release)
		}
	}
	return releases
}

// Release returns the release with the given name
func (c *Cluster) Release(name string) *Release {
	for _, chart := range c.charts {
		for _, release := range chart.Releases() {
			if release.Name() == name {
				return release
			}
		}
	}
	return nil
}
