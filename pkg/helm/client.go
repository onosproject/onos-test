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

package helm

import (
	"github.com/onosproject/onos-test/pkg/kubernetes"
)

var clients = make(map[string]Client)

// Namespace returns a client for the given namespace
func Namespace(namespace ...string) Client {
	var ns string
	if len(namespace) == 0 {
		ns = kubernetes.GetNamespaceFromEnv()
	} else {
		ns = namespace[0]
	}

	client, ok := clients[ns]
	if !ok {
		client = &helmClient{
			Client: kubernetes.Namespace(ns),
			charts: make(map[string]*Chart),
		}
		clients[ns] = client
	}
	return client
}

// Client is a Helm client
type Client interface {
	ChartClient
	ReleaseClient
	kubernetes.Client
}

// ChartClient is a Helm chart client
type ChartClient interface {
	// Charts returns a list of charts in the namespace
	Charts() []*Chart

	// Chart gets a chart in the namespace
	Chart(name string) *Chart
}

// ReleaseClient is a Helm release client
type ReleaseClient interface {
	// Releases returns a list of releases in the namespace
	Releases() []*Release

	// Release gets a chart release in the namespace
	Release(name string) *Release
}

// helmClient is an implementation of the Client interface
type helmClient struct {
	kubernetes.Client
	charts map[string]*Chart
}

// Charts returns a list of charts in the cluster
func (c *helmClient) Charts() []*Chart {
	charts := make([]*Chart, 0, len(c.charts))
	for _, chart := range c.charts {
		charts = append(charts, chart)
	}
	return charts
}

// Chart returns a chart
func (c *helmClient) Chart(name string) *Chart {
	chart, ok := c.charts[name]
	if !ok {
		chart = newChart(name, c.Client)
		c.charts[name] = chart
	}
	return chart
}

// Releases returns a list of releases
func (c *helmClient) Releases() []*Release {
	releases := make([]*Release, 0)
	for _, chart := range c.charts {
		for _, release := range chart.Releases() {
			releases = append(releases, release)
		}
	}
	return releases
}

// Release returns the release with the given name
func (c *helmClient) Release(name string) *Release {
	for _, chart := range c.charts {
		for _, release := range chart.Releases() {
			if release.Name() == name {
				return release
			}
		}
	}
	return nil
}
