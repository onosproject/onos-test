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
	"github.com/atomix/atomix-go-client/pkg/client"
	"github.com/onosproject/onos-test/pkg/onit/cluster"
)

// Partitions is an Atomix partition group
type Partitions interface {
	// List returns a list of partitions in the group
	List() []Partition

	// Get returns the environment for a partition service by name
	Get(name string) Partition

	// Connect connects to the partition group
	Connect() (*client.PartitionGroup, error)
}

// clusterPartitions is an implementation of the Partitions interface
type clusterPartitions struct {
	partitions *cluster.Partitions
}

func (e *clusterPartitions) List() []Partition {
	clusterPartitions := e.partitions.Partitions()
	partitions := make([]Partition, len(clusterPartitions))
	for i, partition := range clusterPartitions {
		partitions[i] = &clusterPartition{
			clusterService: &clusterService{
				service: partition.Service,
			},
		}
	}
	return partitions
}

func (e *clusterPartitions) Get(name string) Partition {
	return &clusterPartition{
		clusterService: &clusterService{
			service: e.partitions.Partition(name).Service,
		},
	}
}

func (e *clusterPartitions) Connect() (*client.PartitionGroup, error) {
	return e.partitions.Connect()
}
