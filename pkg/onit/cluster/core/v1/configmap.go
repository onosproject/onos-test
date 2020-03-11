package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var ConfigMapKind = clustermetav1.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "ConfigMap",
}

var ConfigMapResource = clustermetav1.Resource{
	Kind: ConfigMapKind,
	Name: "ConfigMap",
	ObjectFactory: func() runtime.Object {
		return &corev1.ConfigMap{}
	},
	ObjectsFactory: func() runtime.Object {
		return &corev1.ConfigMapList{}
	},
}

type ConfigMapsClient interface {
	Get(name string) (*ConfigMap, error)
	List() ([]*ConfigMap, error)
}

// newConfigMapsClient creates a new ConfigMapsClient
func newConfigMapsClient(objects clustermetav1.ObjectsClient) ConfigMapsClient {
	return &configMapsClient{
		ObjectsClient: objects,
	}
}

// configMapsClient implements the ConfigMapsClient interface
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

// newConfigMap creates a new ConfigMap resource
func newConfigMap(object *clustermetav1.Object) *ConfigMap {
	return &ConfigMap{
		Object:     object,
		Data:       object.Object.(*corev1.ConfigMap).Data,
		BinaryData: object.Object.(*corev1.ConfigMap).BinaryData,
	}
}

// ConfigMap provides functions for querying a config map
type ConfigMap struct {
	*clustermetav1.Object
	Data       map[string]string
	BinaryData map[string][]byte
}
