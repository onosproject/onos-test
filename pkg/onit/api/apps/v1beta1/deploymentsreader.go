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
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/api/meta/v1"
)

type DeploymentsReader interface {
	Get(name string) (*Deployment, error)
	List() ([]*Deployment, error)
}

func NewDeploymentsReader(objects clustermetav1.ObjectsClient) DeploymentsReader {
	return &deploymentsReader{
		ObjectsClient: objects,
	}
}

type deploymentsReader struct {
	clustermetav1.ObjectsClient
}

func (c *deploymentsReader) Get(name string) (*Deployment, error) {
	object, err := c.ObjectsClient.Get(name, DeploymentResource)
	if err != nil {
		return nil, err
	}
	return NewDeployment(object), nil
}

func (c *deploymentsReader) List() ([]*Deployment, error) {
	objects, err := c.ObjectsClient.List(DeploymentResource)
	if err != nil {
		return nil, err
	}
	deployments := make([]*Deployment, len(objects))
	for i, object := range objects {
		deployments[i] = NewDeployment(object)
	}
	return deployments, nil
}
