package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type ServicesClient interface {
	Get(name string) (*Service, error)
	List() ([]*Service, error)
}

func newServicesClient(objects clustermetav1.ObjectsClient) ServicesClient {
	return &servicesClient{
		ObjectsClient: objects,
	}
}

type servicesClient struct {
	clustermetav1.ObjectsClient
}

func (c *servicesClient) Get(name string) (*Service, error) {
	object, err := c.ObjectsClient.Get(name, ServiceResource)
	if err != nil {
		return nil, err
	}
	return newService(object), nil
}

func (c *servicesClient) List() ([]*Service, error) {
	objects, err := c.ObjectsClient.List(ServiceResource)
	if err != nil {
		return nil, err
	}
	services := make([]*Service, len(objects))
	for i, object := range objects {
		services[i] = newService(object)
	}
	return services, nil
}
