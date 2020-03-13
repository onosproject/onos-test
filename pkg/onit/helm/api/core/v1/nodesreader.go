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
	"github.com/onosproject/onos-test/pkg/onit/helm/api/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

type NodesReader interface {
	Get(name string) (*Node, error)
	List() ([]*Node, error)
}

func NewNodesReader(client resource.Client, filter resource.Filter) NodesReader {
	return &nodesReader{
		Client: client,
		filter: filter,
	}
}

type nodesReader struct {
	resource.Client
	filter resource.Filter
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
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   NodeKind.Group,
			Version: NodeKind.Version,
			Kind:    NodeKind.Kind,
		}, node.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    NodeKind.Group,
				Resource: NodeResource.Name,
			}, name)
		}
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

	results := make([]*Node, 0, len(list.Items))
	for _, node := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   NodeKind.Group,
			Version: NodeKind.Version,
			Kind:    NodeKind.Kind,
		}, node.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			results = append(results, NewNode(&node, c.Client))
		}
	}
	return results, nil
}
