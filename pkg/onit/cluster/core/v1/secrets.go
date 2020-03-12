package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type SecretsClient interface {
	Get(name string) (*Secret, error)
	List() ([]*Secret, error)
}

func newSecretsClient(objects clustermetav1.ObjectsClient) SecretsClient {
	return &secretsClient{
		ObjectsClient: objects,
	}
}

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
