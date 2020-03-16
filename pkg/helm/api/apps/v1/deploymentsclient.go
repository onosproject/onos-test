// Code generated by onit-generate. DO NOT EDIT.

package v1

import (
	"github.com/onosproject/onos-test/pkg/helm/api/resource"
)

type DeploymentsClient interface {
	Deployments() DeploymentsReader
}

func NewDeploymentsClient(resources resource.Client, filter resource.Filter) DeploymentsClient {
	return &deploymentsClient{
		Client: resources,
		filter: filter,
	}
}

type deploymentsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *deploymentsClient) Deployments() DeploymentsReader {
	return NewDeploymentsReader(c.Client, c.filter)
}
