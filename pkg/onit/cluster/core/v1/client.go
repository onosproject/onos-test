package v1

import (
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
	ConfigMaps() ConfigMapsClient
	Endpoints() EndpointsClient
	Nodes() NodesClient
	Pods() PodsClient
	Secrets() SecretsClient
	Services() ServicesClient
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}

func (c *client) ConfigMaps() ConfigMapsClient {
	return newConfigMapsClient(c.ObjectsClient)
}

func (c *client) Endpoints() EndpointsClient {
	return newEndpointsClient(c.ObjectsClient)
}

func (c *client) Nodes() NodesClient {
	return newNodesClient(c.ObjectsClient)
}

func (c *client) Pods() PodsClient {
	return newPodsClient(c.ObjectsClient)
}

func (c *client) Secrets() SecretsClient {
	return newSecretsClient(c.ObjectsClient)
}

func (c *client) Services() ServicesClient {
	return newServicesClient(c.ObjectsClient)
}
