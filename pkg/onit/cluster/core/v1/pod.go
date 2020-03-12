package v1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var PodKind = clustermetav1.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Pod",
}

var PodResource = clustermetav1.Resource{
	Kind: PodKind,
	Name: "Pod",
	ObjectFactory: func() runtime.Object {
		return &corev1.Pod{}
	},
	ObjectsFactory: func() runtime.Object {
		return &corev1.PodList{}
	},
}

func newPod(object *clustermetav1.Object) *Pod {
	return &Pod{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*corev1.Pod).Spec,
	}
}

type Pod struct {
	*clustercorev1.PodSet
	Spec corev1.PodSpec
}
