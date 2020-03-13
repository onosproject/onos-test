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
	"github.com/onosproject/onos-test/pkg/helm/api/resource"
)

type Client interface {
	ConfigMapsClient
	EndpointsClient
	NodesClient
	PodsClient
	SecretsClient
	ServicesClient
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client:           resources,
		ConfigMapsClient: NewConfigMapsClient(resources, filter),
		EndpointsClient:  NewEndpointsClient(resources, filter),
		NodesClient:      NewNodesClient(resources, filter),
		PodsClient:       NewPodsClient(resources, filter),
		SecretsClient:    NewSecretsClient(resources, filter),
		ServicesClient:   NewServicesClient(resources, filter),
	}
}

type client struct {
	resource.Client
	ConfigMapsClient
	EndpointsClient
	NodesClient
	PodsClient
	SecretsClient
	ServicesClient
}
