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

package batch

import (
	batchv1 "github.com/onosproject/onos-test/pkg/onit/api/batch/v1"
	batchv1beta1 "github.com/onosproject/onos-test/pkg/onit/api/batch/v1beta1"
	batchv2alpha1 "github.com/onosproject/onos-test/pkg/onit/api/batch/v2alpha1"
	metav1 "github.com/onosproject/onos-test/pkg/onit/api/meta/v1"
)

type Client interface {
	BatchV1() batchv1.Client
	BatchV1Beta1() batchv1beta1.Client
	BatchV2Alpha1() batchv2alpha1.Client
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}

func (c *client) BatchV1() batchv1.Client {
	return batchv1.NewClient(c.ObjectsClient)
}

func (c *client) BatchV1Beta1() batchv1beta1.Client {
	return batchv1beta1.NewClient(c.ObjectsClient)
}

func (c *client) BatchV2Alpha1() batchv2alpha1.Client {
	return batchv2alpha1.NewClient(c.ObjectsClient)
}
