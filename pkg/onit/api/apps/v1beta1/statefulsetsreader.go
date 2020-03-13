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

package v1beta1

import (
	"github.com/onosproject/onos-test/pkg/onit/api/resource"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type StatefulSetsReader interface {
	Get(name string) (*StatefulSet, error)
	List() ([]*StatefulSet, error)
}

func NewStatefulSetsReader(client resource.Client) StatefulSetsReader {
	return &statefulSetsReader{
		Client: client,
	}
}

type statefulSetsReader struct {
	resource.Client
}

func (c *statefulSetsReader) Get(name string) (*StatefulSet, error) {
	statefulSet := &appsv1beta1.StatefulSet{}
	err := c.Clientset().
		AppsV1beta1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(StatefulSetResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(statefulSet)
	if err != nil {
		return nil, err
	}
	return NewStatefulSet(statefulSet, c.Client), nil
}

func (c *statefulSetsReader) List() ([]*StatefulSet, error) {
	list := &appsv1beta1.StatefulSetList{}
	err := c.Clientset().
		AppsV1beta1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(StatefulSetResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*StatefulSet, len(list.Items))
	for i, statefulSet := range list.Items {
		results[i] = NewStatefulSet(&statefulSet, c.Client)
	}
	return results, nil
}
