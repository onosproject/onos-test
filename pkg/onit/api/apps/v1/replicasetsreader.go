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

type ReplicaSetsReader interface {
	Get(name string) (*ReplicaSet, error)
	List() ([]*ReplicaSet, error)
}

func NewReplicaSetsReader(client resource.Client) ReplicaSetsReader {
	return &replicaSetsReader{
		Client: client,
	}
}

type replicaSetsReader struct {
	resource.Client
}

func (c *replicaSetsReader) Get(name string) (*ReplicaSet, error) {
	replicaSet := &appsv1.ReplicaSet{}
	err := c.Clientset().
		AppsV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(ReplicaSetResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(replicaSet)
	if err != nil {
		return nil, err
	}
	return NewReplicaSet(replicaSet, c.Client), nil
}

func (c *replicaSetsReader) List() ([]*ReplicaSet, error) {
	list := &appsv1.ReplicaSetList{}
	err := c.Clientset().
		AppsV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(ReplicaSetResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*ReplicaSet, len(list.Items))
	for i, replicaSet := range list.Items {
		results[i] = NewReplicaSet(&replicaSet, c.Client)
	}
	return results, nil
}
