package v1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
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
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*appsv1.DaemonSet).Spec,
	}
}

type DaemonSet struct {
	*clustercorev1.PodSet
	Spec appsv1.DaemonSetSpec
}
