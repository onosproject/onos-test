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
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type SecretsReader interface {
	Get(name string) (*Secret, error)
	List() ([]*Secret, error)
}

func NewSecretsReader(objects clustermetav1.ObjectsClient) SecretsReader {
	return &secretsReader{
		ObjectsClient: objects,
	}
}

type secretsReader struct {
	clustermetav1.ObjectsClient
}

func (c *secretsReader) Get(name string) (*Secret, error) {
	object, err := c.ObjectsClient.Get(name, SecretResource)
	if err != nil {
		return nil, err
	}
	return NewSecret(object), nil
}

func (c *secretsReader) List() ([]*Secret, error) {
	objects, err := c.ObjectsClient.List(SecretResource)
	if err != nil {
		return nil, err
	}
	secrets := make([]*Secret, len(objects))
	for i, object := range objects {
		secrets[i] = NewSecret(object)
	}
	return secrets, nil
}
