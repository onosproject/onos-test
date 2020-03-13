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

type appsClient apps.Client
type batchClient batch.Client
type coreClient core.Client
type extensionsClient extensions.Client
type networkingClient networking.Client

type Client interface {
	resource.Client
	appsClient
	batchClient
	coreClient
	extensionsClient
	networkingClient
}

func NewClient(resources resource.Client) Client {
	return &client{
		Client:           resources,
		appsClient:       apps.NewClient(resources),
		batchClient:      batch.NewClient(resources),
		coreClient:       core.NewClient(resources),
		extensionsClient: extensions.NewClient(resources),
		networkingClient: networking.NewClient(resources),
	}
}

type client struct {
	resource.Client
	appsClient
	batchClient
	coreClient
	extensionsClient
	networkingClient
}
