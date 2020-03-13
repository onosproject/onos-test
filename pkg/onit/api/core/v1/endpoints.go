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
	corev1 "k8s.io/api/core/v1"
)

var EndpointsKind = resource.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Endpoints",
}

var EndpointsResource = resource.Type{
	Kind: EndpointsKind,
	Name: "endpoints",
}

func NewEndpoints(endpoints *corev1.Endpoints, client resource.Client) *Endpoints {
	return &Endpoints{
		Resource:  resource.NewResource(endpoints.ObjectMeta, EndpointsKind, client),
		Endpoints: endpoints,
	}
}

type Endpoints struct {
	*resource.Resource
	Endpoints *corev1.Endpoints
}
