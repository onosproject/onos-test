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

type DaemonSetsClient interface {
	Get(name string) (*DaemonSet, error)
	List() ([]*DaemonSet, error)
}

// newDaemonSetsClient creates a new DaemonSetsClient
func newDaemonSetsClient(objects clustermetav1.ObjectsClient) DaemonSetsClient {
	return &daemonSetsClient{
		ObjectsClient: objects,
	}
}

// daemonSetsClient implements the DaemonSetsClient interface
type daemonSetsClient struct {
	clustermetav1.ObjectsClient
}

func (c *daemonSetsClient) Get(name string) (*DaemonSet, error) {
	object, err := c.ObjectsClient.Get(name, DaemonSetResource)
	if err != nil {
		return nil, err
	}
	return newDaemonSet(object), nil
}

func (c *daemonSetsClient) List() ([]*DaemonSet, error) {
	objects, err := c.ObjectsClient.List(DaemonSetResource)
	if err != nil {
		return nil, err
	}
	daemonSets := make([]*DaemonSet, len(objects))
	for i, object := range objects {
		daemonSets[i] = newDaemonSet(object)
	}
	return daemonSets, nil
}

// newDaemonSet creates a new DaemonSet resource
func newDaemonSet(object *clustermetav1.Object) *DaemonSet {
	return &DaemonSet{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*appsv1.DaemonSet).Spec,
	}
}

// DaemonSet provides functions for querying a daemon set
type DaemonSet struct {
	*clustercorev1.PodSet
	Spec appsv1.DaemonSetSpec
}
