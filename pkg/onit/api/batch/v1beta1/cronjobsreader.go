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

package v1beta1

import (
	"github.com/onosproject/onos-test/pkg/onit/api/resource"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type CronJobsReader interface {
	Get(name string) (*CronJob, error)
	List() ([]*CronJob, error)
}

func NewCronJobsReader(client resource.Client) CronJobsReader {
	return &cronJobsReader{
		Client: client,
	}
}

type cronJobsReader struct {
	resource.Client
}

func (c *cronJobsReader) Get(name string) (*CronJob, error) {
	cronJob := &batchv1beta1.CronJob{}
	err := c.Clientset().
		BatchV1beta1().
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
	}
	return NewCronJob(cronJob, c.Client), nil
}

func (c *cronJobsReader) List() ([]*CronJob, error) {
	list := &batchv1beta1.CronJobList{}
	err := c.Clientset().
		BatchV1beta1().
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

	results := make([]*CronJob, len(list.Items))
	for i, cronJob := range list.Items {
		results[i] = NewCronJob(&cronJob, c.Client)
	}
	return results, nil
}
