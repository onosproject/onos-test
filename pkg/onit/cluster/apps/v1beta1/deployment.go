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
	appsv1 "github.com/onosproject/onos-test/pkg/onit/cluster/apps/v1"
	corev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var DeploymentKind = clustermetav1.Kind{
	Group:   "apps",
	Version: "v1beta1",
	Kind:    "Deployment",
}

var DeploymentResource = clustermetav1.Resource{
	Kind: DeploymentKind,
	Name: "Deployment",
	ObjectFactory: func() runtime.Object {
		return &appsv1beta1.Deployment{}
	},
	ObjectsFactory: func() runtime.Object {
		return &appsv1beta1.DeploymentList{}
	},
}

func NewDeployment(object *clustermetav1.Object) *Deployment {
	return &Deployment{
		Object:            object,
		Deployment:        object.Object.(*appsv1beta1.Deployment),
		ReplicaSetsClient: appsv1.NewReplicaSetsClient(object),
		PodsClient:        corev1.NewPodsClient(object),
	}
}

type Deployment struct {
	*clustermetav1.Object
	Deployment *appsv1beta1.Deployment
	appsv1.ReplicaSetsClient
	corev1.PodsClient
}
