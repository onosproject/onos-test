package v1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var ReplicaSetKind = clustermetav1.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "ReplicaSet",
}

var ReplicaSetResource = clustermetav1.Resource{
	Kind: ReplicaSetKind,
	Name: "ReplicaSet",
	ObjectFactory: func() runtime.Object {
		return &appsv1.ReplicaSet{}
	},
	ObjectsFactory: func() runtime.Object {
		return &appsv1.ReplicaSetList{}
	},
}

func newReplicaSet(object *clustermetav1.Object) *ReplicaSet {
	return &ReplicaSet{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*appsv1.ReplicaSet).Spec,
	}
}

type ReplicaSet struct {
	*clustercorev1.PodSet
	Spec appsv1.ReplicaSetSpec
}
