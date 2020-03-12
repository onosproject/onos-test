package v1beta1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type IngressesClient interface {
	Get(name string) (*Ingress, error)
	List() ([]*Ingress, error)
}

func newIngressesClient(objects clustermetav1.ObjectsClient) IngressesClient {
	return &ingressesClient{
		ObjectsClient: objects,
	}
}

type ingressesClient struct {
	clustermetav1.ObjectsClient
}

func (c *ingressesClient) Get(name string) (*Ingress, error) {
	object, err := c.ObjectsClient.Get(name, IngressResource)
	if err != nil {
		return nil, err
	}
	return newIngress(object), nil
}

func (c *ingressesClient) List() ([]*Ingress, error) {
	objects, err := c.ObjectsClient.List(IngressResource)
	if err != nil {
		return nil, err
	}
	ingresses := make([]*Ingress, len(objects))
	for i, object := range objects {
		ingresses[i] = newIngress(object)
	}
	return ingresses, nil
}
