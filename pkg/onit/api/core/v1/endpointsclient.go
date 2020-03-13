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
)

type EndpointsClient interface {
	Endpoints() EndpointsReader
}

func NewEndpointsClient(resources resource.Client, filter resource.Filter) EndpointsClient {
	return &endpointsClient{
		Client: resources,
		filter: filter,
	}
}

type endpointsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *endpointsClient) Endpoints() EndpointsReader {
	return NewEndpointsReader(c.Client, c.filter)
}