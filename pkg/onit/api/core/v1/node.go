// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/api/meta/v1"
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

func NewNode(object *clustermetav1.Object) *Node {
	return &Node{
		Object: object,
		Node:   object.Object.(*corev1.Node),
	}
}

type Node struct {
	*clustermetav1.Object
	Node *corev1.Node
}
