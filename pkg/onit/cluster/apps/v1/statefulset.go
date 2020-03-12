package v1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var StatefulSetKind = clustermetav1.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "StatefulSet",
}

var StatefulSetResource = clustermetav1.Resource{
	Kind: StatefulSetKind,
	Name: "StatefulSet",
	ObjectFactory: func() runtime.Object {
		return &appsv1.StatefulSet{}
	},
	ObjectsFactory: func() runtime.Object {
		return &appsv1.StatefulSetList{}
	},
}

func newStatefulSet(object *clustermetav1.Object) *StatefulSet {
	return &StatefulSet{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*appsv1.StatefulSet).Spec,
	}
}

type StatefulSet struct {
	*clustercorev1.PodSet
	Spec appsv1.StatefulSetSpec
}
