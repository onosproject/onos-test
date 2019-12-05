// Copyright 2019-present Open Networking Foundation.
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

package kube

import (
	"errors"
	"github.com/onosproject/onos-test/pkg/util/random"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NamespaceEnv is the environment variable for setting the k8s namespace
const NamespaceEnv = "NAMESPACE"

// GetAPI returns the Kubernetes API for the given namespace
func GetAPI(namespace string) (API, error) {
	config, err := GetRestConfig()
	if err != nil {
		return nil, err
	}
	client, err := client.New(config, client.Options{})
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &kubeAPI{
		namespace: namespace,
		config:    config,
		client:    client,
		clientset: clientset,
	}, nil
}

// GetAPIOrDie returns the Kubernetes API for the given namespace and panics if no configuration is found
func GetAPIOrDie(namespace string) API {
	if api, err := GetAPI(namespace); err != nil {
		panic(err)
	} else {
		return api
	}
}

// GetAPIFromEnv returns the Kubernetes API for the current environment
func GetAPIFromEnv() (API, error) {
	namespace := os.Getenv(NamespaceEnv)
	if namespace == "" {
		namespace = random.NewPetName(2)
	}
	return GetAPI(namespace)
}

// GetAPIFromEnvOrDie returns the Kubernetes API for the current environment and panics if no API could be found
func GetAPIFromEnvOrDie() API {
	if api, err := GetAPIFromEnv(); err != nil {
		panic(err)
	} else {
		return api
	}
}

// APIProvider is an interface for types to provide the Kubernetes API
type APIProvider interface {
	// API returns the API
	API() API
}

// API exposes the Kubernetes API to tests
type API interface {
	// Namespace returns the Kubernetes namespace
	Namespace() string

	// Config returns the Kubernetes REST configuration
	Config() *rest.Config

	// Client returns the Kubernetes controller runtime client
	Client() client.Client

	// Clientset returns the Kubernetes Go clientset
	Clientset() *kubernetes.Clientset
}

// kubeAPI provides the Kubernetes API
type kubeAPI struct {
	namespace string
	config    *rest.Config
	client    client.Client
	clientset *kubernetes.Clientset
}

func (k *kubeAPI) Namespace() string {
	return k.namespace
}

func (k *kubeAPI) Config() *rest.Config {
	return k.config
}

func (k *kubeAPI) Client() client.Client {
	return k.client
}

func (k *kubeAPI) Clientset() *kubernetes.Clientset {
	return k.clientset
}

// GetClient returns a Kubernetes client
func GetClient() (client.Client, error) {
	config, err := GetRestConfig()
	if err != nil {
		return nil, err
	}
	return client.New(config, client.Options{})
}

// GetClientset returns a Kubernetes clientset
func GetClientset() (*kubernetes.Clientset, error) {
	config, err := GetRestConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
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
