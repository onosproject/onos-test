package v1beta1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var IngressKind = clustermetav1.Kind{
	Group:   "extensions",
	Version: "v1beta1",
	Kind:    "Ingress",
}

var IngressResource = clustermetav1.Resource{
	Kind: IngressKind,
	Name: "Ingress",
	ObjectFactory: func() runtime.Object {
		return &extensionsv1beta1.Ingress{}
	},
	ObjectsFactory: func() runtime.Object {
		return &extensionsv1beta1.IngressList{}
	},
}

func newIngress(object *clustermetav1.Object) *Ingress {
	return &Ingress{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*extensionsv1beta1.Ingress).Spec,
	}
}

type Ingress struct {
	*clustercorev1.PodSet
	Spec extensionsv1beta1.IngressSpec
}
