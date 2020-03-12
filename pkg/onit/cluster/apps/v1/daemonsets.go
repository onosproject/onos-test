package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type DaemonSetsClient interface {
	Get(name string) (*DaemonSet, error)
	List() ([]*DaemonSet, error)
}

func newDaemonSetsClient(objects clustermetav1.ObjectsClient) DaemonSetsClient {
	return &daemonSetsClient{
		ObjectsClient: objects,
	}
}

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
