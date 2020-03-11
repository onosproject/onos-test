package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var EndpointsKind = clustermetav1.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Endpoints",
}

var EndpointsResource = clustermetav1.Resource{
	Kind: EndpointsKind,
	Name: "Endpoints",
	ObjectFactory: func() runtime.Object {
		return &corev1.Endpoints{}
	},
	ObjectsFactory: func() runtime.Object {
		return &corev1.EndpointsList{}
	},
}

type EndpointsClient interface {
	Get(name string) (*Endpoints, error)
	List() ([]*Endpoints, error)
}

// newEndpointsClient creates a new EndpointsClient
func newEndpointsClient(objects clustermetav1.ObjectsClient) EndpointsClient {
	return &endpointsClient{
		ObjectsClient: objects,
	}
}

// endpointsClient implements the EndpointsClient interface
type endpointsClient struct {
	clustermetav1.ObjectsClient
}

func (c *endpointsClient) Get(name string) (*Endpoints, error) {
	object, err := c.ObjectsClient.Get(name, EndpointsResource)
	if err != nil {
		return nil, err
	}
	return newEndpoints(object), nil
}

func (c *endpointsClient) List() ([]*Endpoints, error) {
	objects, err := c.ObjectsClient.List(EndpointsResource)
	if err != nil {
		return nil, err
	}
	endpoints := make([]*Endpoints, len(objects))
	for i, object := range objects {
		endpoints[i] = newEndpoints(object)
	}
	return endpoints, nil
}

// newEndpoints creates a new Endpoints resource
func newEndpoints(object *clustermetav1.Object) *Endpoints {
	return &Endpoints{
		Object:  object,
		Subsets: object.Object.(*corev1.Endpoints).Subsets,
	}
}

// Endpoints provides functions for querying endpoints
type Endpoints struct {
	*clustermetav1.Object
	Subsets []corev1.EndpointSubset
}
