package v1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var SecretKind = clustermetav1.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Secret",
}

var SecretResource = clustermetav1.Resource{
	Kind: SecretKind,
	Name: "Secret",
	ObjectFactory: func() runtime.Object {
		return &corev1.Secret{}
	},
	ObjectsFactory: func() runtime.Object {
		return &corev1.SecretList{}
	},
}

func newSecret(object *clustermetav1.Object) *Secret {
	return &Secret{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*corev1.Secret).Spec,
	}
}

type Secret struct {
	*clustercorev1.PodSet
	Spec corev1.SecretSpec
}
