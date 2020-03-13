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

type ServicesReader interface {
	Get(name string) (*Service, error)
	List() ([]*Service, error)
}

func NewServicesReader(client resource.Client, filter resource.Filter) ServicesReader {
	return &servicesReader{
		Client: client,
		filter: filter,
	}
}

type servicesReader struct {
	resource.Client
	filter resource.Filter
}

func (c *servicesReader) Get(name string) (*Service, error) {
	service := &corev1.Service{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(ServiceResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(service)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ServiceKind.Group,
			Version: ServiceKind.Version,
			Kind:    ServiceKind.Kind,
		}, service.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    ServiceKind.Group,
				Resource: ServiceResource.Name,
			}, name)
		}
	}
	return NewService(service, c.Client), nil
}

func (c *servicesReader) List() ([]*Service, error) {
	list := &corev1.ServiceList{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(ServiceResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Service, 0, len(list.Items))
	for _, service := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ServiceKind.Group,
			Version: ServiceKind.Version,
			Kind:    ServiceKind.Kind,
		}, service.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			results = append(results, NewService(&service, c.Client))
		}
	}
	return results, nil
}
