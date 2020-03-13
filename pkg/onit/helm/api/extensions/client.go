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

package extensions

import (
	extensionsv1beta1 "github.com/onosproject/onos-test/pkg/onit/helm/api/extensions/v1beta1"
	"github.com/onosproject/onos-test/pkg/onit/helm/api/resource"
)

type Client interface {
	V1beta1() extensionsv1beta1.Client
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client: resources,
		filter: filter,
	}
}

type client struct {
	resource.Client
	filter resource.Filter
}

func (c *client) V1beta1() extensionsv1beta1.Client {
	return extensionsv1beta1.NewClient(c.Client, c.filter)
}
