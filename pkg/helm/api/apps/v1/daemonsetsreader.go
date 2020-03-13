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
	"github.com/onosproject/onos-test/pkg/helm/api/resource"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

type DaemonSetsReader interface {
	Get(name string) (*DaemonSet, error)
	List() ([]*DaemonSet, error)
}

func NewDaemonSetsReader(client resource.Client, filter resource.Filter) DaemonSetsReader {
	return &daemonSetsReader{
		Client: client,
		filter: filter,
	}
}

type daemonSetsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *daemonSetsReader) Get(name string) (*DaemonSet, error) {
	daemonSet := &appsv1.DaemonSet{}
	err := c.Clientset().
		AppsV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(DaemonSetResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(daemonSet)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   DaemonSetKind.Group,
			Version: DaemonSetKind.Version,
			Kind:    DaemonSetKind.Kind,
		}, daemonSet.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    DaemonSetKind.Group,
				Resource: DaemonSetResource.Name,
			}, name)
		}
	}
	return NewDaemonSet(daemonSet, c.Client), nil
}

func (c *daemonSetsReader) List() ([]*DaemonSet, error) {
	list := &appsv1.DaemonSetList{}
	err := c.Clientset().
		AppsV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(DaemonSetResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*DaemonSet, 0, len(list.Items))
	for _, daemonSet := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   DaemonSetKind.Group,
			Version: DaemonSetKind.Version,
			Kind:    DaemonSetKind.Kind,
		}, daemonSet.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			results = append(results, NewDaemonSet(&daemonSet, c.Client))
		}
	}
	return results, nil
}
