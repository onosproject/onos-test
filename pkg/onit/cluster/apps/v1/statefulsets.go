package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type StatefulSetsClient interface {
	Get(name string) (*StatefulSet, error)
	List() ([]*StatefulSet, error)
}

func newStatefulSetsClient(objects clustermetav1.ObjectsClient) StatefulSetsClient {
	return &statefulSetsClient{
		ObjectsClient: objects,
	}
}

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
