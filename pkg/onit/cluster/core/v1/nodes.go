package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type NodesClient interface {
	Get(name string) (*Node, error)
	List() ([]*Node, error)
}

func newNodesClient(objects clustermetav1.ObjectsClient) NodesClient {
	return &nodesClient{
		ObjectsClient: objects,
	}
}

type nodesClient struct {
	clustermetav1.ObjectsClient
}

func (c *nodesClient) Get(name string) (*Node, error) {
	object, err := c.ObjectsClient.Get(name, NodeResource)
	if err != nil {
		return nil, err
	}
	return newNode(object), nil
}

func (c *nodesClient) List() ([]*Node, error) {
	objects, err := c.ObjectsClient.List(NodeResource)
	if err != nil {
		return nil, err
	}
	nodes := make([]*Node, len(objects))
	for i, object := range objects {
		nodes[i] = newNode(object)
	}
	return nodes, nil
}
