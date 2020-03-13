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
	"github.com/onosproject/onos-test/pkg/onit/api/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type NodesReader interface {
	Get(name string) (*Node, error)
	List() ([]*Node, error)
}

func NewNodesReader(client resource.Client) NodesReader {
	return &nodesReader{
		Client: client,
	}
}

type nodesReader struct {
	resource.Client
}

func (c *nodesReader) Get(name string) (*Node, error) {
	node := &corev1.Node{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(NodeResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(node)
	if err != nil {
		return nil, err
	}
	return NewNode(node, c.Client), nil
}

func (c *nodesReader) List() ([]*Node, error) {
	list := &corev1.NodeList{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(NodeResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Node, len(list.Items))
	for i, node := range list.Items {
		results[i] = NewNode(&node, c.Client)
	}
	return results, nil
}
