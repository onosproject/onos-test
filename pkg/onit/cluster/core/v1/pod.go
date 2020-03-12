package v1

import (
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
		Object: object,
		Pod: object.Object.(*corev1.Pod),
	}
}

type Pod struct {
	*clustermetav1.Object
	Pod *corev1.Pod
}
