package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type ReplicaSetsClient interface {
	Get(name string) (*ReplicaSet, error)
	List() ([]*ReplicaSet, error)
}

func newReplicaSetsClient(objects clustermetav1.ObjectsClient) ReplicaSetsClient {
	return &replicaSetsClient{
		ObjectsClient: objects,
	}
}

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
