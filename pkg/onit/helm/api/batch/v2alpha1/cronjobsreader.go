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

package v2alpha1

import (
	"github.com/onosproject/onos-test/pkg/onit/helm/api/resource"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

type CronJobsReader interface {
	Get(name string) (*CronJob, error)
	List() ([]*CronJob, error)
}

func NewCronJobsReader(client resource.Client, filter resource.Filter) CronJobsReader {
	return &cronJobsReader{
		Client: client,
		filter: filter,
	}
}

type cronJobsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *cronJobsReader) Get(name string) (*CronJob, error) {
	cronJob := &batchv2alpha1.CronJob{}
	err := c.Clientset().
		BatchV2alpha1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(CronJobResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(cronJob)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   CronJobKind.Group,
			Version: CronJobKind.Version,
			Kind:    CronJobKind.Kind,
		}, cronJob.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    CronJobKind.Group,
				Resource: CronJobResource.Name,
			}, name)
		}
	}
	return NewCronJob(cronJob, c.Client), nil
}

func (c *cronJobsReader) List() ([]*CronJob, error) {
	list := &batchv2alpha1.CronJobList{}
	err := c.Clientset().
		BatchV2alpha1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(CronJobResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*CronJob, 0, len(list.Items))
	for _, cronJob := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   CronJobKind.Group,
			Version: CronJobKind.Version,
			Kind:    CronJobKind.Kind,
		}, cronJob.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			results = append(results, NewCronJob(&cronJob, c.Client))
		}
	}
	return results, nil
}
