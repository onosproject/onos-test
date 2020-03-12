package v1beta1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var StatefulSetKind = clustermetav1.Kind{
	Group:   "apps",
	Version: "v1beta1",
	Kind:    "StatefulSet",
}

var StatefulSetResource = clustermetav1.Resource{
	Kind: StatefulSetKind,
	Name: "StatefulSet",
	ObjectFactory: func() runtime.Object {
		return &appsv1beta1.StatefulSet{}
	},
	ObjectsFactory: func() runtime.Object {
		return &appsv1beta1.StatefulSetList{}
	},
}

func newStatefulSet(object *clustermetav1.Object) *StatefulSet {
	return &StatefulSet{
		Object: object,
		StatefulSet: object.Object.(*appsv1beta1.StatefulSet),
	}
}

type StatefulSet struct {
	*clustermetav1.Object
	StatefulSet *appsv1beta1.StatefulSet
}