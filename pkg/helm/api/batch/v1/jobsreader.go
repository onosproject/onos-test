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
	"github.com/onosproject/onos-test/pkg/helm/api/resource"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

type JobsReader interface {
	Get(name string) (*Job, error)
	List() ([]*Job, error)
}

func NewJobsReader(client resource.Client, filter resource.Filter) JobsReader {
	return &jobsReader{
		Client: client,
		filter: filter,
	}
}

type jobsReader struct {
	resource.Client
	filter resource.Filter
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
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   JobKind.Group,
			Version: JobKind.Version,
			Kind:    JobKind.Kind,
		}, job.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    JobKind.Group,
				Resource: JobResource.Name,
			}, name)
		}
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

	results := make([]*Job, 0, len(list.Items))
	for _, job := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   JobKind.Group,
			Version: JobKind.Version,
			Kind:    JobKind.Kind,
		}, job.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			results = append(results, NewJob(&job, c.Client))
		}
	}
	return results, nil
}
