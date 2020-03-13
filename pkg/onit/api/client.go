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

package api

import (
	apps "github.com/onosproject/onos-test/pkg/onit/api/apps"
	batch "github.com/onosproject/onos-test/pkg/onit/api/batch"
	core "github.com/onosproject/onos-test/pkg/onit/api/core"
	extensions "github.com/onosproject/onos-test/pkg/onit/api/extensions"
	networking "github.com/onosproject/onos-test/pkg/onit/api/networking"
	"github.com/onosproject/onos-test/pkg/onit/api/resource"
)

type Client interface {
	resource.Client
	Apps() apps.Client
	Batch() batch.Client
	Core() core.Client
	Extensions() extensions.Client
	Networking() networking.Client
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

func (c *client) Apps() apps.Client {
	return apps.NewClient(c.Client, c.filter)
}

func (c *client) Batch() batch.Client {
	return batch.NewClient(c.Client, c.filter)
}

func (c *client) Core() core.Client {
	return core.NewClient(c.Client, c.filter)
}

func (c *client) Extensions() extensions.Client {
	return extensions.NewClient(c.Client, c.filter)
}

func (c *client) Networking() networking.Client {
	return networking.NewClient(c.Client, c.filter)
}
