package kubetest

import (
	"github.com/onosproject/onos-test/pkg/new/util/k8s"
	"k8s.io/client-go/rest"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getKubeAPI returns the Kubernetes API for the current environment
func getKubeAPI() KubeAPI {
	namespace := os.Getenv(testNamespaceEnv)
	config, err := k8s.GetRestConfig()
	if err != nil {
		panic(err)
	}
	client, err := client.New(config, client.Options{})
	if err != nil {
		panic(err)
	}
	return &kubeAPI{
		namespace: namespace,
		config:    config,
		client:    client,
	}
}

// KubeAPIProvider is an interface for types to provide the Kubernetes API
type KubeAPIProvider interface {
	// KubeAPI returns the KubeAPI
	KubeAPI() KubeAPI
}

// KubeAPI exposes the Kubernetes API to tests
type KubeAPI interface {
	// Namespace returns the Kubernetes namespace
	Namespace() string

	// Config returns the Kubernetes REST configuration
	Config() *rest.Config

	// Client returns the Kubernetes controller runtime client
	Client() client.Client
}

// kubeAPI provides the Kubernetes API
type kubeAPI struct {
	namespace string
	config    *rest.Config
	client    client.Client
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