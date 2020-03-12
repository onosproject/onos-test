package cluster

import (
	"github.com/onosproject/onos-test/pkg/kube"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

func newChart(name string, parent metav1.ObjectsClient) *Chart {
	return &Chart{
		API:      parent,
		name:     name,
		releases: make(map[string]*Release),
	}
}

// Chart is a Helm chart
type Chart struct {
	kube.API
	name       string
	repository string
	releases   map[string]*Release
}

// Name returns the chart name
func (c *Chart) Name() string {
	return c.name
}

// SetRepository sets the chart's repository URL
func (c *Chart) SetRepository(url string) {
	c.repository = url
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
