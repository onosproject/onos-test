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

func newEndpoints(object *clustermetav1.Object) *Endpoints {
	return &Endpoints{
		Object: object,
		Endpoints: object.Object.(*corev1.Endpoints),
	}
}

type Endpoints struct {
	*clustermetav1.Object
	Endpoints *corev1.Endpoints
}
