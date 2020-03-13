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
	corev1 "github.com/onosproject/onos-test/pkg/onit/api/core/v1"
	"github.com/onosproject/onos-test/pkg/onit/api/resource"
	appsv1 "k8s.io/api/apps/v1"
)

var ReplicaSetKind = resource.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "ReplicaSet",
}

var ReplicaSetResource = resource.Type{
	Kind: ReplicaSetKind,
	Name: "replicasets",
}

func NewReplicaSet(replicaSet *appsv1.ReplicaSet, client resource.Client) *ReplicaSet {
	return &ReplicaSet{
		Resource:   resource.NewResource(replicaSet.ObjectMeta, ReplicaSetKind, client),
		ReplicaSet: replicaSet,
		PodsClient: corev1.NewPodsClient(client, resource.NewUIDFilter(replicaSet.UID)),
	}
}

type ReplicaSet struct {
	*resource.Resource
	ReplicaSet *appsv1.ReplicaSet
	corev1.PodsClient
}
