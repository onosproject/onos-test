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
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type JobsReader interface {
	Get(name string) (*Job, error)
	List() ([]*Job, error)
}

func NewJobsReader(client resource.Client) JobsReader {
	return &jobsReader{
		Client: client,
	}
}

type jobsReader struct {
	resource.Client
}

func (c *jobsReader) Get(name string) (*Job, error) {
	job := &batchv1.Job{}
	err := c.Clientset().
		BatchV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(JobResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(job)
	if err != nil {
		return nil, err
	}
	return NewJob(job, c.Client), nil
}

func (c *jobsReader) List() ([]*Job, error) {
	list := &batchv1.JobList{}
	err := c.Clientset().
		BatchV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(JobResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Job, len(list.Items))
	for i, job := range list.Items {
		results[i] = NewJob(&job, c.Client)
	}
	return results, nil
}
