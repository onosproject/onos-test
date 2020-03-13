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
	"errors"
	"github.com/onosproject/onos-test/pkg/util/random"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// NamespaceEnv is the environment variable for setting the k8s namespace
const NamespaceEnv = "POD_NAMESPACE"

// GetNamespaceFromEnv gets the Kubernetes namespace from the environment
func GetNamespaceFromEnv() string {
	namespace := os.Getenv(NamespaceEnv)
	if namespace == "" {
		namespace = random.NewPetName(2)
	}
	return namespace
}

// Namespace returns the Helm namespace
func Namespace(namespace ...string) Client {
	config := GetRestConfigOrDie()
	ns := GetNamespaceFromEnv()
	if len(namespace) > 0 {
		ns = namespace[0]
	}
	return &kubernetesClient{
		namespace: ns,
		client:    kubernetes.NewForConfigOrDie(config),
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
	namespace string
	client    *kubernetes.Clientset
}

func (c *kubernetesClient) Namespace() string {
	return c.namespace
}

func (c *kubernetesClient) Clientset() *kubernetes.Clientset {
	return c.client
}

// GetRestConfigOrDie returns the Kubernetes REST API configuration
func GetRestConfigOrDie() *rest.Config {
	config, err := GetRestConfig()
	if err != nil {
		panic(err)
	}
	return config
}

// GetRestConfig returns the Kubernetes REST API configuration
func GetRestConfig() (*rest.Config, error) {
	restconfig, err := rest.InClusterConfig()
	if err == nil {
		return restconfig, nil
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home := getHomeDir()
		if home == "" {
			return nil, errors.New("no home directory configured")
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// use the current context in kubeconfig
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

// getHomeDir returns the user's home directory if defined by environment variables
func getHomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
