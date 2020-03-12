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

func newService(object *clustermetav1.Object) *Service {
	return &Service{
		Object: object,
		Service: object.Object.(*corev1.Service),
	}
}

type Service struct {
	*clustermetav1.Object
	Service *corev1.Service
}
