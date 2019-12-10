// Copyright 2019-present Open Networking Foundation.
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

package cluster

import "strings"

const (
	partitionType  = "partition"
	groupLabel     = "group"
	partitionLabel = "partition"
)

func newPartition(cluster *Cluster, name string) *Partition {
	labels := getLabels(partitionType)
	labels[groupLabel] = name[:strings.LastIndex(name, "-")]
	labels[partitionLabel] = name[strings.LastIndex(name, "-")+1:]
	deployment := newDeployment(cluster)
	deployment.SetName(name)
	deployment.SetLabels(labels)

	return &Partition{
		Deployment: deployment,
	}
}

// Partition provides methods for querying a database partition
type Partition struct {
	*Deployment
}
