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

type ReplicaSetsClient interface {
	Get(name string) (*ReplicaSet, error)
	List() ([]*ReplicaSet, error)
}

// newReplicaSetsClient creates a new ReplicaSetsClient
func newReplicaSetsClient(objects clustermetav1.ObjectsClient) ReplicaSetsClient {
	return &replicaSetsClient{
		ObjectsClient: objects,
	}
}

// replicaSetsClient implements the ReplicaSetsClient interface
type replicaSetsClient struct {
	clustermetav1.ObjectsClient
}

func (c *replicaSetsClient) Get(name string) (*ReplicaSet, error) {
	object, err := c.ObjectsClient.Get(name, ReplicaSetResource)
	if err != nil {
		return nil, err
	}
	return newReplicaSet(object), nil
}

func (c *replicaSetsClient) List() ([]*ReplicaSet, error) {
	objects, err := c.ObjectsClient.List(ReplicaSetResource)
	if err != nil {
		return nil, err
	}
	replicaSets := make([]*ReplicaSet, len(objects))
	for i, object := range objects {
		replicaSets[i] = newReplicaSet(object)
	}
	return replicaSets, nil
}

// newReplicaSet creates a new ReplicaSet resource
func newReplicaSet(object *clustermetav1.Object) *ReplicaSet {
	return &ReplicaSet{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*appsv1.ReplicaSet).Spec,
	}
}

// ReplicaSet provides functions for querying a replica set
type ReplicaSet struct {
	*clustercorev1.PodSet
	Spec appsv1.ReplicaSetSpec
}
