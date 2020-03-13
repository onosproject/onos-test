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
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/api/meta/v1"
)

type ServicesReader interface {
	Get(name string) (*Service, error)
	List() ([]*Service, error)
}

func NewServicesReader(objects clustermetav1.ObjectsClient) ServicesReader {
	return &servicesReader{
		ObjectsClient: objects,
	}
}

type servicesReader struct {
	clustermetav1.ObjectsClient
}

func (c *servicesReader) Get(name string) (*Service, error) {
	object, err := c.ObjectsClient.Get(name, ServiceResource)
	if err != nil {
		return nil, err
	}
	return NewService(object), nil
}

func (c *servicesReader) List() ([]*Service, error) {
	objects, err := c.ObjectsClient.List(ServiceResource)
	if err != nil {
		return nil, err
	}
	services := make([]*Service, len(objects))
	for i, object := range objects {
		services[i] = NewService(object)
	}
	return services, nil
}
