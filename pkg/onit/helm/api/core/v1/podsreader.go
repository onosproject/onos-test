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
	"github.com/onosproject/onos-test/pkg/onit/helm/api/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

type PodsReader interface {
	Get(name string) (*Pod, error)
	List() ([]*Pod, error)
}

func NewPodsReader(client resource.Client, filter resource.Filter) PodsReader {
	return &podsReader{
		Client: client,
		filter: filter,
	}
}

type podsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *podsReader) Get(name string) (*Pod, error) {
	pod := &corev1.Pod{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(PodResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(pod)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   PodKind.Group,
			Version: PodKind.Version,
			Kind:    PodKind.Kind,
		}, pod.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    PodKind.Group,
				Resource: PodResource.Name,
			}, name)
		}
	}
	return NewPod(pod, c.Client), nil
}

func (c *podsReader) List() ([]*Pod, error) {
	list := &corev1.PodList{}
	err := c.Clientset().
		CoreV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(PodResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Pod, 0, len(list.Items))
	for _, pod := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   PodKind.Group,
			Version: PodKind.Version,
			Kind:    PodKind.Kind,
		}, pod.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			results = append(results, NewPod(&pod, c.Client))
		}
	}
	return results, nil
}
