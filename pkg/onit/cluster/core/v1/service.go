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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var ServiceKind = clustermetav1.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Service",
}

var ServiceResource = clustermetav1.Resource{
	Kind: ServiceKind,
	Name: "Service",
	ObjectFactory: func() runtime.Object {
		return &corev1.Service{}
	},
	ObjectsFactory: func() runtime.Object {
		return &corev1.ServiceList{}
	},
}

func NewService(object *clustermetav1.Object) *Service {
	return &Service{
		Object:          object,
		Service:         object.Object.(*corev1.Service),
		EndpointsClient: NewEndpointsClient(object),
	}
}

type Service struct {
	*clustermetav1.Object
	Service *corev1.Service
	EndpointsClient
}
