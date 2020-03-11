package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var NodeKind = clustermetav1.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Node",
}

var NodeResource = clustermetav1.Resource{
	Kind: NodeKind,
	Name: "Node",
	ObjectFactory: func() runtime.Object {
		return &corev1.Node{}
	},
	ObjectsFactory: func() runtime.Object {
		return &corev1.NodeList{}
	},
}

type NodesClient interface {
	Get(name string) (*Node, error)
	List() ([]*Node, error)
}

// newNodesClient creates a new NodesClient
func newNodesClient(objects clustermetav1.ObjectsClient) NodesClient {
	return &nodesClient{
		ObjectsClient: objects,
	}
}

// nodesClient implements the NodesClient interface
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

// newNode creates a new Node resource
func newNode(object *clustermetav1.Object) *Node {
	return &Node{
		Object: object,
		Spec:   object.Object.(*corev1.Node).Spec,
	}
}

// Node provides functions for querying a node
type Node struct {
	*clustermetav1.Object
	Spec corev1.NodeSpec
}
