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
	metav1 "github.com/onosproject/onos-test/pkg/onit/api/meta/v1"
)

type Client interface {
	ConfigMapsClient
	EndpointsClient
	NodesClient
	PodsClient
	SecretsClient
	ServicesClient
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient:    objects,
		ConfigMapsClient: NewConfigMapsClient(objects),

		EndpointsClient: NewEndpointsClient(objects),

		NodesClient: NewNodesClient(objects),

		PodsClient: NewPodsClient(objects),

		SecretsClient: NewSecretsClient(objects),

		ServicesClient: NewServicesClient(objects),
	}
}

type client struct {
	metav1.ObjectsClient
	ConfigMapsClient
	EndpointsClient
	NodesClient
	PodsClient
	SecretsClient
	ServicesClient
}
