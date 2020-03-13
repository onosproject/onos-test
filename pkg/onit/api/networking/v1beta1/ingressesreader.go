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
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type IngressesReader interface {
	Get(name string) (*Ingress, error)
	List() ([]*Ingress, error)
}

func NewIngressesReader(client resource.Client) IngressesReader {
	return &ingressesReader{
		Client: client,
	}
}

type ingressesReader struct {
	resource.Client
}

func (c *ingressesReader) Get(name string) (*Ingress, error) {
	ingress := &networkingv1beta1.Ingress{}
	err := c.Clientset().
		NetworkingV1beta1().
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
	}
	return NewIngress(ingress, c.Client), nil
}

func (c *ingressesReader) List() ([]*Ingress, error) {
	list := &networkingv1beta1.IngressList{}
	err := c.Clientset().
		NetworkingV1beta1().
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

	results := make([]*Ingress, len(list.Items))
	for i, ingress := range list.Items {
		results[i] = NewIngress(&ingress, c.Client)
	}
	return results, nil
}
