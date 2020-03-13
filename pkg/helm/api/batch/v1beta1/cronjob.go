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
	"github.com/onosproject/onos-test/pkg/helm/api/resource"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
)

var CronJobKind = resource.Kind{
	Group:   "batch",
	Version: "v1beta1",
	Kind:    "CronJob",
}

var CronJobResource = resource.Type{
	Kind: CronJobKind,
	Name: "cronjobs",
}

func NewCronJob(cronJob *batchv1beta1.CronJob, client resource.Client) *CronJob {
	return &CronJob{
		Resource: resource.NewResource(cronJob.ObjectMeta, CronJobKind, client),
		CronJob:  cronJob,
	}
}

type CronJob struct {
	*resource.Resource
	CronJob *batchv1beta1.CronJob
}
