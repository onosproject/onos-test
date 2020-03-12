package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type PodsClient interface {
	Get(name string) (*Pod, error)
	List() ([]*Pod, error)
}

func newPodsClient(objects clustermetav1.ObjectsClient) PodsClient {
	return &podsClient{
		ObjectsClient: objects,
	}
}

type podsClient struct {
	clustermetav1.ObjectsClient
}

func (c *podsClient) Get(name string) (*Pod, error) {
	object, err := c.ObjectsClient.Get(name, PodResource)
	if err != nil {
		return nil, err
	}
	return newPod(object), nil
}

func (c *podsClient) List() ([]*Pod, error) {
	objects, err := c.ObjectsClient.List(PodResource)
	if err != nil {
		return nil, err
	}
	pods := make([]*Pod, len(objects))
	for i, object := range objects {
		pods[i] = newPod(object)
	}
	return pods, nil
}
