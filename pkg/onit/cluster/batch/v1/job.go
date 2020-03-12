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
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var JobKind = clustermetav1.Kind{
	Group:   "batch",
	Version: "v1",
	Kind:    "Job",
}

var JobResource = clustermetav1.Resource{
	Kind: JobKind,
	Name: "Job",
	ObjectFactory: func() runtime.Object {
		return &batchv1.Job{}
	},
	ObjectsFactory: func() runtime.Object {
		return &batchv1.JobList{}
	},
}

func NewJob(object *clustermetav1.Object) *Job {
	return &Job{
		Object: object,
		Job: object.Object.(*batchv1.Job),
	}
}

type Job struct {
	*clustermetav1.Object
	Job *batchv1.Job
}
