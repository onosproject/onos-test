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
	appsv1 "github.com/onosproject/onos-test/pkg/onit/api/apps/v1"
	corev1 "github.com/onosproject/onos-test/pkg/onit/api/core/v1"
	"github.com/onosproject/onos-test/pkg/onit/api/resource"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
)

var DeploymentKind = resource.Kind{
	Group:   "apps",
	Version: "v1beta1",
	Kind:    "Deployment",
}

var DeploymentResource = resource.Type{
	Kind: DeploymentKind,
	Name: "deployments",
}

func NewDeployment(deployment *appsv1beta1.Deployment, client resource.Client) *Deployment {
	return &Deployment{
		Resource:          resource.NewResource(deployment.ObjectMeta, DeploymentKind, client),
		Deployment:        deployment,
		ReplicaSetsClient: appsv1.NewReplicaSetsClient(client, resource.NewUIDFilter(deployment.UID)),
		PodsClient:        corev1.NewPodsClient(client, resource.NewUIDFilter(deployment.UID)),
	}
}

type Deployment struct {
	*resource.Resource
	Deployment *appsv1beta1.Deployment
	appsv1.ReplicaSetsClient
	corev1.PodsClient
}
