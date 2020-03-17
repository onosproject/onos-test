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
	"helm.sh/helm/v3/pkg/action"
	"log"
)

var clients = make(map[string]Client)

// Namespace returns the Helm namespace
func Namespace() string {
	return kubernetes.GetNamespaceFromEnv()
}

// Helm returns the Helm client
func Helm() Client {
	return getClient(kubernetes.GetNamespaceFromEnv())
}

// getClient returns the client for the given namespace
func getClient(namespace string) Client {
	client, ok := clients[namespace]
	if !ok {
		config, err := getConfig(namespace)
		if err != nil {
			panic(err)
		}
		client = &helmClient{
			Client: kubernetes.NewClient(namespace),
			charts: make(map[string]*Chart),
			config: config,
		}
		clients[namespace] = client
	}
	return client
}

// getConfig gets the Helm configuration for the given namespace
func getConfig(namespace string) (*action.Configuration, error) {
	config := &action.Configuration{}
	if err := config.Init(settings.RESTClientGetter(), namespace, "memory", log.Printf); err != nil {
		return nil, err
	}
	return config, nil
}

// Client is a Helm client
type Client interface {
	ChartClient
	ReleaseClient

	// Namespace returns the client for the given namespace
	Namespace(namespace string) Client
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
	config *action.Configuration
}

func (c *helmClient) Namespace(namespace string) Client {
	return getClient(namespace)
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
		chart = newChart(name, c.Client, c.config)
		c.charts[name] = chart
	}
	return chart
}

// Releases returns a list of releases
func (c *helmClient) Releases() []*Release {
	releases := make([]*Release, 0)
	for _, chart := range c.charts {
		releases = append(releases, chart.Releases()...)
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
