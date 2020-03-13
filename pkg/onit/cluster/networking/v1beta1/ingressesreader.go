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
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type IngressesReader interface {
	Get(name string) (*Ingress, error)
	List() ([]*Ingress, error)
}

func NewIngressesReader(objects clustermetav1.ObjectsClient) IngressesReader {
	return &ingressesReader{
		ObjectsClient: objects,
	}
}

type ingressesReader struct {
	clustermetav1.ObjectsClient
}

func (c *ingressesReader) Get(name string) (*Ingress, error) {
	object, err := c.ObjectsClient.Get(name, IngressResource)
	if err != nil {
		return nil, err
	}
	return NewIngress(object), nil
}

func (c *ingressesReader) List() ([]*Ingress, error) {
	objects, err := c.ObjectsClient.List(IngressResource)
	if err != nil {
		return nil, err
	}
	ingresses := make([]*Ingress, len(objects))
	for i, object := range objects {
		ingresses[i] = NewIngress(object)
	}
	return ingresses, nil
}
