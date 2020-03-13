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

package cluster

import (
	apps "github.com/onosproject/onos-test/pkg/onit/cluster/apps"
	batch "github.com/onosproject/onos-test/pkg/onit/cluster/batch"
	core "github.com/onosproject/onos-test/pkg/onit/cluster/core"
	extensions "github.com/onosproject/onos-test/pkg/onit/cluster/extensions"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	networking "github.com/onosproject/onos-test/pkg/onit/cluster/networking"
)

type appsClient apps.Client
type batchClient batch.Client
type coreClient core.Client
type extensionsClient extensions.Client
type networkingClient networking.Client

type Client interface {
	appsClient
	batchClient
	coreClient
	extensionsClient
	networkingClient
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient:    objects,
		appsClient:       apps.NewClient(objects),
		batchClient:      batch.NewClient(objects),
		coreClient:       core.NewClient(objects),
		extensionsClient: extensions.NewClient(objects),
		networkingClient: networking.NewClient(objects),
	}
}

type client struct {
	metav1.ObjectsClient
	appsClient
	batchClient
	coreClient
	extensionsClient
	networkingClient
}
