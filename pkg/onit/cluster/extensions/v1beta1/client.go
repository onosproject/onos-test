package v1beta1

import (
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
    Ingresses() IngressesClient
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}


func (c *client) Ingresses() IngressesClient {
	return newIngressesClient(c.ObjectsClient)
}

