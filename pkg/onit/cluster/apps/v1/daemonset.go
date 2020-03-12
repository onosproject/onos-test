package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var DaemonSetKind = clustermetav1.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "DaemonSet",
}

var DaemonSetResource = clustermetav1.Resource{
	Kind: DaemonSetKind,
	Name: "DaemonSet",
	ObjectFactory: func() runtime.Object {
		return &appsv1.DaemonSet{}
	},
	ObjectsFactory: func() runtime.Object {
		return &appsv1.DaemonSetList{}
	},
}

func newDaemonSet(object *clustermetav1.Object) *DaemonSet {
	return &DaemonSet{
		Object: object,
		DaemonSet: object.Object.(*appsv1.DaemonSet),
	}
}

type DaemonSet struct {
	*clustermetav1.Object
	DaemonSet *appsv1.DaemonSet
}
