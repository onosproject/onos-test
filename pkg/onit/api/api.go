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
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/onit/api/resource"
)

// NewFromEnv gets the cluster from the environment
func NewFromEnv() *API {
	return New(kube.GetAPIFromEnvOrDie())
}

// New returns a new onit Env
func New(api kube.API) *API {
	return &API{
		Client: NewClient(api, resource.NoFilter),
		charts: make(map[string]*Chart),
	}
}

var _ Client = &API{}

// API facilitates modifying subsystems in Kubernetes
type API struct {
	Client
	charts map[string]*Chart
}

// Charts returns a list of charts in the cluster
func (c *API) Charts() []*Chart {
	charts := make([]*Chart, 0, len(c.charts))
	for _, chart := range c.charts {
		charts = append(charts, chart)
	}
	return charts
}

// Chart returns a chart
func (c *API) Chart(name string) *Chart {
	chart, ok := c.charts[name]
	if !ok {
		chart = newChart(name, c.Client)
		c.charts[name] = chart
	}
	return chart
}

// Releases returns a list of releases
func (c *API) Releases() []*Release {
	releases := make([]*Release, 0)
	for _, chart := range c.charts {
		for _, release := range chart.Releases() {
			releases = append(releases, release)
		}
	}
	return releases
}

// Release returns the release with the given name
func (c *API) Release(name string) *Release {
	for _, chart := range c.charts {
		for _, release := range chart.Releases() {
			if release.Name() == name {
				return release
			}
		}
	}
	return nil
}
