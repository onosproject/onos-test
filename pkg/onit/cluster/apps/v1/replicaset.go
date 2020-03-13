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
	corev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var ReplicaSetKind = clustermetav1.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "ReplicaSet",
}

var ReplicaSetResource = clustermetav1.Resource{
	Kind: ReplicaSetKind,
	Name: "ReplicaSet",
	ObjectFactory: func() runtime.Object {
		return &appsv1.ReplicaSet{}
	},
	ObjectsFactory: func() runtime.Object {
		return &appsv1.ReplicaSetList{}
	},
}

func NewReplicaSet(object *clustermetav1.Object) *ReplicaSet {
	return &ReplicaSet{
		Object:     object,
		ReplicaSet: object.Object.(*appsv1.ReplicaSet),
		PodsClient: corev1.NewPodsClient(object),
	}
}

type ReplicaSet struct {
	*clustermetav1.Object
	ReplicaSet *appsv1.ReplicaSet
	corev1.PodsClient
}
