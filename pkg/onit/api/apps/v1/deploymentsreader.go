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
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type DeploymentsReader interface {
	Get(name string) (*Deployment, error)
	List() ([]*Deployment, error)
}

func NewDeploymentsReader(client resource.Client) DeploymentsReader {
	return &deploymentsReader{
		Client: client,
	}
}

type deploymentsReader struct {
	resource.Client
}

func (c *deploymentsReader) Get(name string) (*Deployment, error) {
	deployment := &appsv1.Deployment{}
	err := c.Clientset().
		AppsV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(DeploymentResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(deployment)
	if err != nil {
		return nil, err
	}
	return NewDeployment(deployment, c.Client), nil
}

func (c *deploymentsReader) List() ([]*Deployment, error) {
	list := &appsv1.DeploymentList{}
	err := c.Clientset().
		AppsV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(DeploymentResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Deployment, len(list.Items))
	for i, deployment := range list.Items {
		results[i] = NewDeployment(&deployment, c.Client)
	}
	return results, nil
}
