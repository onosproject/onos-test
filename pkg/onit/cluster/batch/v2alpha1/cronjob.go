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
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

var CronJobKind = clustermetav1.Kind{
	Group:   "batch",
	Version: "v2alpha1",
	Kind:    "CronJob",
}

var CronJobResource = clustermetav1.Resource{
	Kind: CronJobKind,
	Name: "CronJob",
	ObjectFactory: func() runtime.Object {
		return &batchv2alpha1.CronJob{}
	},
	ObjectsFactory: func() runtime.Object {
		return &batchv2alpha1.CronJobList{}
	},
}

func NewCronJob(object *clustermetav1.Object) *CronJob {
	return &CronJob{
		Object:  object,
		CronJob: object.Object.(*batchv2alpha1.CronJob),
	}
}

type CronJob struct {
	*clustermetav1.Object
	CronJob *batchv2alpha1.CronJob
}
