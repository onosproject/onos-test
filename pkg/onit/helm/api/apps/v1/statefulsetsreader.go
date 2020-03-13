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
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

type StatefulSetsReader interface {
	Get(name string) (*StatefulSet, error)
	List() ([]*StatefulSet, error)
}

func NewStatefulSetsReader(client resource.Client, filter resource.Filter) StatefulSetsReader {
	return &statefulSetsReader{
		Client: client,
		filter: filter,
	}
}

type statefulSetsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *statefulSetsReader) Get(name string) (*StatefulSet, error) {
	statefulSet := &appsv1.StatefulSet{}
	err := c.Clientset().
		AppsV1().
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
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   StatefulSetKind.Group,
			Version: StatefulSetKind.Version,
			Kind:    StatefulSetKind.Kind,
		}, statefulSet.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    StatefulSetKind.Group,
				Resource: StatefulSetResource.Name,
			}, name)
		}
	}
	return NewStatefulSet(statefulSet, c.Client), nil
}

func (c *statefulSetsReader) List() ([]*StatefulSet, error) {
	list := &appsv1.StatefulSetList{}
	err := c.Clientset().
		AppsV1().
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

	results := make([]*StatefulSet, 0, len(list.Items))
	for _, statefulSet := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   StatefulSetKind.Group,
			Version: StatefulSetKind.Version,
			Kind:    StatefulSetKind.Kind,
		}, statefulSet.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			results = append(results, NewStatefulSet(&statefulSet, c.Client))
		}
	}
	return results, nil
}
