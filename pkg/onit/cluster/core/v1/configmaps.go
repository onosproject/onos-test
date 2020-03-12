package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type ConfigMapsClient interface {
	Get(name string) (*ConfigMap, error)
	List() ([]*ConfigMap, error)
}

func newConfigMapsClient(objects clustermetav1.ObjectsClient) ConfigMapsClient {
	return &configMapsClient{
		ObjectsClient: objects,
	}
}

type configMapsClient struct {
	clustermetav1.ObjectsClient
}

func (c *configMapsClient) Get(name string) (*ConfigMap, error) {
	object, err := c.ObjectsClient.Get(name, ConfigMapResource)
	if err != nil {
		return nil, err
	}
	return newConfigMap(object), nil
}

func (c *configMapsClient) List() ([]*ConfigMap, error) {
	objects, err := c.ObjectsClient.List(ConfigMapResource)
	if err != nil {
		return nil, err
	}
	configMaps := make([]*ConfigMap, len(objects))
	for i, object := range objects {
		configMaps[i] = newConfigMap(object)
	}
	return configMaps, nil
}
