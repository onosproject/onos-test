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
	"github.com/onosproject/onos-test/pkg/helm/api/resource"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var IngressKind = resource.Kind{
	Group:   "networking",
	Version: "v1beta1",
	Kind:    "Ingress",
}

var IngressResource = resource.Type{
	Kind: IngressKind,
	Name: "ingresses",
}

func NewIngress(ingress *networkingv1beta1.Ingress, client resource.Client) *Ingress {
	return &Ingress{
		Resource: resource.NewResource(ingress.ObjectMeta, IngressKind, client),
		Object:   ingress,
	}
}

type Ingress struct {
	*resource.Resource
	Object *networkingv1beta1.Ingress
}

func (r *Ingress) Delete() error {
	return r.Clientset().
		NetworkingV1beta1().
		RESTClient().
		Delete().
		Namespace(r.Namespace).
		Resource(IngressResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}
