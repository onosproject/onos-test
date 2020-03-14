// Copyright 2020-present Open Networking Foundation.
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
	"github.com/onosproject/onos-test/pkg/helm/api/resource"
)

func newChart(name string, client resource.Client) *Chart {
	return &Chart{
		Client:   client,
		name:     name,
		releases: make(map[string]*Release),
	}
}

// Chart is a Helm chart
type Chart struct {
	resource.Client
	ReleaseClient
	name       string
	repository string
	releases   map[string]*Release
}

// Name returns the chart name
func (c *Chart) Name() string {
	return c.name
}

// SetRepository sets the chart's repository URL
func (c *Chart) SetRepository(url string) *Chart {
	c.repository = url
	return c
}

// Repository returns the chart's repository URL
func (c *Chart) Repository() string {
	return c.repository
}

// Releases returns a list of releases of the chart
func (c *Chart) Releases() []*Release {
	releases := make([]*Release, 0, len(c.releases))
	for _, release := range c.releases {
		releases = append(releases, release)
	}
	return releases
}

// Release returns the release with the given name
func (c *Chart) Release(name string) *Release {
	release, ok := c.releases[name]
	if !ok {
		release = newRelease(name, c)
		c.releases[name] = release
	}
	return release
}
