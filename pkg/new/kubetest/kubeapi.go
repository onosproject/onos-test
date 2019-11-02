package kubetest

import (
	"github.com/onosproject/onos-test/pkg/new/util/k8s"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getKubeAPI returns the Kubernetes API for the current environment
func getKubeAPI() KubeAPI {
	config, err := k8s.GetRestConfig()
	if err != nil {
		panic(err)
	}
	client, err := client.New(config, client.Options{})
	if err != nil {
		panic(err)
	}
	return &kubeAPI{
		config: config,
		client: client,
	}
}

// KubeAPIProvider is an interface for types to provide the Kubernetes API
type KubeAPIProvider interface {
	// KubeAPI returns the KubeAPI
	KubeAPI() KubeAPI
}

// KubeAPI exposes the Kubernetes API to tests
type KubeAPI interface {
	// Config returns the Kubernetes REST configuration
	Config() *rest.Config

	// Client returns the Kubernetes controller runtime client
	Client() client.Client
}

// kubeAPI provides the Kubernetes API
type kubeAPI struct {
	config *rest.Config
	client client.Client
}

func (k *kubeAPI) Config() *rest.Config {
	return k.config
}

func (k *kubeAPI) Client() client.Client {
	return k.client
}
