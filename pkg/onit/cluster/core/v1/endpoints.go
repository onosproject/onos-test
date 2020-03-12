package v1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
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
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*corev1.Endpoints).Spec,
	}
}

type Endpoints struct {
	*clustercorev1.PodSet
	Spec corev1.EndpointsSpec
}
