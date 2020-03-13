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
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

type IngressesReader interface {
	Get(name string) (*Ingress, error)
	List() ([]*Ingress, error)
}

func NewIngressesReader(client resource.Client, filter resource.Filter) IngressesReader {
	return &ingressesReader{
		Client: client,
		filter: filter,
	}
}

type ingressesReader struct {
	resource.Client
	filter resource.Filter
}

func (c *ingressesReader) Get(name string) (*Ingress, error) {
	ingress := &extensionsv1beta1.Ingress{}
	err := c.Clientset().
		ExtensionsV1beta1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(IngressResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(ingress)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   IngressKind.Group,
			Version: IngressKind.Version,
			Kind:    IngressKind.Kind,
		}, ingress.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    IngressKind.Group,
				Resource: IngressResource.Name,
			}, name)
		}
	}
	return NewIngress(ingress, c.Client), nil
}

func (c *ingressesReader) List() ([]*Ingress, error) {
	list := &extensionsv1beta1.IngressList{}
	err := c.Clientset().
		ExtensionsV1beta1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(IngressResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Ingress, 0, len(list.Items))
	for _, ingress := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   IngressKind.Group,
			Version: IngressKind.Version,
			Kind:    IngressKind.Kind,
		}, ingress.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			results = append(results, NewIngress(&ingress, c.Client))
		}
	}
	return results, nil
}
