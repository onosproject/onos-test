package v1beta1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
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

type StatefulSetsClient interface {
	Get(name string) (*StatefulSet, error)
	List() ([]*StatefulSet, error)
}

// newStatefulSetsClient creates a new StatefulSetsClient
func newStatefulSetsClient(objects clustermetav1.ObjectsClient) StatefulSetsClient {
	return &statefulSetsClient{
		ObjectsClient: objects,
	}
}

// statefulSetsClient implements the StatefulSetsClient interface
type statefulSetsClient struct {
	clustermetav1.ObjectsClient
}

func (c *statefulSetsClient) Get(name string) (*StatefulSet, error) {
	object, err := c.ObjectsClient.Get(name, StatefulSetResource)
	if err != nil {
		return nil, err
	}
	return newStatefulSet(object), nil
}

func (c *statefulSetsClient) List() ([]*StatefulSet, error) {
	objects, err := c.ObjectsClient.List(StatefulSetResource)
	if err != nil {
		return nil, err
	}
	statefulSets := make([]*StatefulSet, len(objects))
	for i, object := range objects {
		statefulSets[i] = newStatefulSet(object)
	}
	return statefulSets, nil
}

// newStatefulSet creates a new StatefulSet resource
func newStatefulSet(object *clustermetav1.Object) *StatefulSet {
	return &StatefulSet{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*appsv1beta1.StatefulSet).Spec,
	}
}

// StatefulSet provides functions for querying a stateful set
type StatefulSet struct {
	*clustercorev1.PodSet
	Spec appsv1beta1.StatefulSetSpec
}
