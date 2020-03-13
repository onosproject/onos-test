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

package kubernetes

import (
	"github.com/onosproject/onos-test/pkg/kube"
	"k8s.io/client-go/kubernetes"
)

// Namespace returns the Helm namespace
func Namespace(namespace ...string) Client {
	if len(namespace) == 0 {
		return &kubernetesClient{
			api: kube.GetAPIFromEnvOrDie(),
		}
	}
	return &kubernetesClient{
		api: kube.GetAPIOrDie(namespace[0]),
	}
}

// Client is a Kubernetes client
type Client interface {
	// Namespace returns the client namespace
	Namespace() string

	// Clientset returns the client's Clientset
	Clientset() *kubernetes.Clientset
}

// kubernetesClient is an implementation of the Kubernetes Client interface
type kubernetesClient struct {
	api kube.API
}

func (c *kubernetesClient) Namespace() string {
	return c.api.Namespace()
}

func (c *kubernetesClient) Clientset() *kubernetes.Clientset {
	return c.api.Clientset()
}
