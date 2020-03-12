package v1

import (
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
		Object: object,
		ReplicaSet: object.Object.(*appsv1.ReplicaSet),
	}
}

type ReplicaSet struct {
	*clustermetav1.Object
	ReplicaSet *appsv1.ReplicaSet
}
