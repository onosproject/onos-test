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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

type ConfigMapsReader interface {
	Get(name string) (*ConfigMap, error)
	List() ([]*ConfigMap, error)
}

func NewConfigMapsReader(client resource.Client, filter resource.Filter) ConfigMapsReader {
	return &configMapsReader{
		Client: client,
		filter: filter,
	}
}

type configMapsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *configMapsReader) Get(name string) (*ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(ConfigMapResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(configMap)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ConfigMapKind.Group,
			Version: ConfigMapKind.Version,
			Kind:    ConfigMapKind.Kind,
		}, configMap.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    ConfigMapKind.Group,
				Resource: ConfigMapResource.Name,
			}, name)
		}
	}
	return NewConfigMap(configMap, c.Client), nil
}

func (c *configMapsReader) List() ([]*ConfigMap, error) {
	list := &corev1.ConfigMapList{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(ConfigMapResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*ConfigMap, 0, len(list.Items))
	for _, configMap := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ConfigMapKind.Group,
			Version: ConfigMapKind.Version,
			Kind:    ConfigMapKind.Kind,
		}, configMap.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			results = append(results, NewConfigMap(&configMap, c.Client))
		}
	}
	return results, nil
}
