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
	corev1 "github.com/onosproject/onos-test/pkg/onit/helm/api/core/v1"
	"github.com/onosproject/onos-test/pkg/onit/helm/api/resource"
	appsv1 "k8s.io/api/apps/v1"
)

var StatefulSetKind = resource.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "StatefulSet",
}

var StatefulSetResource = resource.Type{
	Kind: StatefulSetKind,
	Name: "statefulsets",
}

func NewStatefulSet(statefulSet *appsv1.StatefulSet, client resource.Client) *StatefulSet {
	return &StatefulSet{
		Resource:      resource.NewResource(statefulSet.ObjectMeta, StatefulSetKind, client),
		StatefulSet:   statefulSet,
		PodsReference: corev1.NewPodsReference(client, resource.NewUIDFilter(statefulSet.UID)),
	}
}

type StatefulSet struct {
	*resource.Resource
	StatefulSet *appsv1.StatefulSet
	corev1.PodsReference
}
