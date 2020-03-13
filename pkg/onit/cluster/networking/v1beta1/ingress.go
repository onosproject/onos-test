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
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var IngressKind = clustermetav1.Kind{
	Group:   "networking",
	Version: "v1beta1",
	Kind:    "Ingress",
}

var IngressResource = clustermetav1.Resource{
	Kind: IngressKind,
	Name: "Ingress",
	ObjectFactory: func() runtime.Object {
		return &networkingv1beta1.Ingress{}
	},
	ObjectsFactory: func() runtime.Object {
		return &networkingv1beta1.IngressList{}
	},
}

func NewIngress(object *clustermetav1.Object) *Ingress {
	return &Ingress{
		Object:  object,
		Ingress: object.Object.(*networkingv1beta1.Ingress),
	}
}

type Ingress struct {
	*clustermetav1.Object
	Ingress *networkingv1beta1.Ingress
}
