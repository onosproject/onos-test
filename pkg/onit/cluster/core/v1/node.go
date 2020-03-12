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

func newNode(object *clustermetav1.Object) *Node {
	return &Node{
		Object: object,
		Node: object.Object.(*corev1.Node),
	}
}

type Node struct {
	*clustermetav1.Object
	Node *corev1.Node
}
