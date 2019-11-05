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

package env

import (
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
)

// Database provides the database environment
type Database interface {
	// Partitions returns all database partitions
	Partitions(group string) []Partition

	// Partition returns the Partition for the given partition
	Partition(name string) Partition
}

var _ Database = &clusterDatabase{}

// clusterDatabase is an implementation of the Database interface
type clusterDatabase struct {
	group *cluster.Partitions
}

func (e *clusterDatabase) Partitions(group string) []Partition {
	clusterPartitions := e.group.Partitions()
	partitions := make([]Partition, len(clusterPartitions))
	for i, partition := range clusterPartitions {
		partitions[i] = e.Partition(partition.Name())
	}
	return partitions
}

func (e clusterDatabase) Partition(name string) Partition {
	return &clusterPartition{
		clusterService: &clusterService{
			service: e.group.Partition(name).Service,
		},
	}
}
