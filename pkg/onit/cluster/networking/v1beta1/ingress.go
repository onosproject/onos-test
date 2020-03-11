package v1beta1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var IngressKind = clustermetav1.Kind{
	Group:   "networking",
	Version: "v1beta1",
	Kind:    "Ingress",
}

var IngressResource = clustermetav1.Resource{
	Kind: IngressKind,
	Name: "Ingress",
	ObjectFactory: func() runtime.Object {
		return &networkingv1beta1.Ingress{}
	},
	ObjectsFactory: func() runtime.Object {
		return &networkingv1beta1.IngressList{}
	},
}

type IngressesClient interface {
	Get(name string) (*Ingress, error)
	List() ([]*Ingress, error)
}

// newIngressesClient creates a new IngressesClient
func newIngressesClient(objects clustermetav1.ObjectsClient) IngressesClient {
	return &ingressesClient{
		ObjectsClient: objects,
	}
}

// ingressesClient implements the IngressesClient interface
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

// newIngress creates a new Ingress resource
func newIngress(object *clustermetav1.Object) *Ingress {
	return &Ingress{
		Object: object,
		Spec:   object.Object.(*networkingv1beta1.Ingress).Spec,
	}
}

// Ingress provides functions for querying an ingress
type Ingress struct {
	*clustermetav1.Object
	Spec networkingv1beta1.IngressSpec
}
