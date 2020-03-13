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
	corev1 "github.com/onosproject/onos-test/pkg/helm/api/core/v1"
	"github.com/onosproject/onos-test/pkg/helm/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReplicaSetsReference interface {
	ReplicaSets() ReplicaSetsReader
	corev1.PodsReference
}

func NewReplicaSetsReference(resources resource.Client, filter resource.Filter) ReplicaSetsReference {
	var ownerFilter resource.Filter = func(kind metav1.GroupVersionKind, meta metav1.ObjectMeta) (bool, error) {
		list, err := NewReplicaSetsReader(resources, filter).List()
		if err != nil {
			return false, err
		}
		for _, owner := range meta.OwnerReferences {
			for _, replicaSets := range list {
				if replicaSets.ReplicaSet.ObjectMeta.UID == owner.UID {
					return true, nil
				}
			}
		}
		return false, nil
	}
	return &replicaSetsReference{
		Client:        resources,
		filter:        filter,
		PodsReference: corev1.NewPodsReference(resources, ownerFilter),
	}
}

type replicaSetsReference struct {
	resource.Client
	filter resource.Filter
	corev1.PodsReference
}

func (c *replicaSetsReference) ReplicaSets() ReplicaSetsReader {
	return NewReplicaSetsReader(c.Client, c.filter)
}
