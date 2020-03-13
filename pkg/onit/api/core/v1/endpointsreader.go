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

type EndpointsReader interface {
	Get(name string) (*Endpoints, error)
	List() ([]*Endpoints, error)
}

func NewEndpointsReader(client resource.Client) EndpointsReader {
	return &endpointsReader{
		Client: client,
	}
}

type endpointsReader struct {
	resource.Client
}

func (c *endpointsReader) Get(name string) (*Endpoints, error) {
	endpoints := &corev1.Endpoints{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(EndpointsResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(endpoints)
	if err != nil {
		return nil, err
	}
	return NewEndpoints(endpoints, c.Client), nil
}

func (c *endpointsReader) List() ([]*Endpoints, error) {
	list := &corev1.EndpointsList{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(EndpointsResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Endpoints, len(list.Items))
	for i, endpoints := range list.Items {
		results[i] = NewEndpoints(&endpoints, c.Client)
	}
	return results, nil
}
