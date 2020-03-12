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
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type CronJobs interface {
	Get(name string) (*CronJob, error)
	List() ([]*CronJob, error)
}

func NewCronJobs(objects clustermetav1.ObjectsClient) CronJobs {
	return &cronJobs{
		ObjectsClient: objects,
	}
}

type cronJobs struct {
	clustermetav1.ObjectsClient
}

func (c *cronJobs) Get(name string) (*CronJob, error) {
	object, err := c.ObjectsClient.Get(name, CronJobResource)
	if err != nil {
		return nil, err
	}
	return NewCronJob(object), nil
}

func (c *cronJobs) List() ([]*CronJob, error) {
	objects, err := c.ObjectsClient.List(CronJobResource)
	if err != nil {
		return nil, err
	}
	cronJobs := make([]*CronJob, len(objects))
	for i, object := range objects {
		cronJobs[i] = NewCronJob(object)
	}
	return cronJobs, nil
}
