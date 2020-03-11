package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var SecretKind = clustermetav1.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Secret",
}

var SecretResource = clustermetav1.Resource{
	Kind: SecretKind,
	Name: "Secret",
	ObjectFactory: func() runtime.Object {
		return &corev1.Secret{}
	},
	ObjectsFactory: func() runtime.Object {
		return &corev1.SecretList{}
	},
}

type SecretsClient interface {
	Get(name string) (*Secret, error)
	List() ([]*Secret, error)
}

// newSecretsClient creates a new SecretsClient
func newSecretsClient(objects clustermetav1.ObjectsClient) SecretsClient {
	return &secretsClient{
		ObjectsClient: objects,
	}
}

// secretsClient implements the SecretsClient interface
type secretsClient struct {
	clustermetav1.ObjectsClient
}

func (c *secretsClient) Get(name string) (*Secret, error) {
	object, err := c.ObjectsClient.Get(name, SecretResource)
	if err != nil {
		return nil, err
	}
	return newSecret(object), nil
}

func (c *secretsClient) List() ([]*Secret, error) {
	objects, err := c.ObjectsClient.List(SecretResource)
	if err != nil {
		return nil, err
	}
	secrets := make([]*Secret, len(objects))
	for i, object := range objects {
		secrets[i] = newSecret(object)
	}
	return secrets, nil
}

// newSecret creates a new Secret resource
func newSecret(object *clustermetav1.Object) *Secret {
	return &Secret{
		Object: object,
		Data:   object.Object.(*corev1.Secret).Data,
	}
}

// Secret provides functions for querying a secret
type Secret struct {
	*clustermetav1.Object
	Data map[string][]byte
}
