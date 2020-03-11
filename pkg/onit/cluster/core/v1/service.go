package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var ServiceKind = clustermetav1.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Service",
}

var ServiceResource = clustermetav1.Resource{
	Kind: ServiceKind,
	Name: "Service",
	ObjectFactory: func() runtime.Object {
		return &corev1.Service{}
	},
	ObjectsFactory: func() runtime.Object {
		return &corev1.ServiceList{}
	},
}

type ServicesClient interface {
	Get(name string) (*Service, error)
	List() ([]*Service, error)
}

// newServicesClient creates a new ServicesClient
func newServicesClient(objects clustermetav1.ObjectsClient) ServicesClient {
	return &servicesClient{
		ObjectsClient: objects,
	}
}

// servicesClient implements the ServicesClient interface
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

// newService creates a new Service resource
func newService(object *clustermetav1.Object) *Service {
	return &Service{
		Object: object,
		Spec:   object.Object.(*corev1.Service).Spec,
	}
}

// Service provides functions for querying a service
type Service struct {
	*clustermetav1.Object
	Spec corev1.ServiceSpec
}
