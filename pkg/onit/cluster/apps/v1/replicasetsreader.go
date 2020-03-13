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
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type ReplicaSetsReader interface {
	Get(name string) (*ReplicaSet, error)
	List() ([]*ReplicaSet, error)
}

func NewReplicaSetsReader(objects clustermetav1.ObjectsClient) ReplicaSetsReader {
	return &replicaSetsReader{
		ObjectsClient: objects,
	}
}

type replicaSetsReader struct {
	clustermetav1.ObjectsClient
}

func (c *replicaSetsReader) Get(name string) (*ReplicaSet, error) {
	object, err := c.ObjectsClient.Get(name, ReplicaSetResource)
	if err != nil {
		return nil, err
	}
	return NewReplicaSet(object), nil
}

func (c *replicaSetsReader) List() ([]*ReplicaSet, error) {
	objects, err := c.ObjectsClient.List(ReplicaSetResource)
	if err != nil {
		return nil, err
	}
	replicaSets := make([]*ReplicaSet, len(objects))
	for i, object := range objects {
		replicaSets[i] = NewReplicaSet(object)
	}
	return replicaSets, nil
}
