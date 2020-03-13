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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type SecretsReader interface {
	Get(name string) (*Secret, error)
	List() ([]*Secret, error)
}

func NewSecretsReader(client resource.Client) SecretsReader {
	return &secretsReader{
		Client: client,
	}
}

type secretsReader struct {
	resource.Client
}

func (c *secretsReader) Get(name string) (*Secret, error) {
	secret := &corev1.Secret{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(SecretResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(secret)
	if err != nil {
		return nil, err
	}
	return NewSecret(secret, c.Client), nil
}

func (c *secretsReader) List() ([]*Secret, error) {
	list := &corev1.SecretList{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(SecretResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Secret, len(list.Items))
	for i, secret := range list.Items {
		results[i] = NewSecret(&secret, c.Client)
	}
	return results, nil
}
