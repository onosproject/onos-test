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

// PartitionsEnv is an Atomix partition group
type PartitionsEnv interface {
	// List returns a list of partitions in the group
	List() []PartitionEnv

	// Get returns the environment for a partition service by name
	Get(name string) PartitionEnv

	// Connect connects to the partition group
	Connect() (*client.PartitionGroup, error)
}

// clusterPartitionsEnv is an implementation of the Partitions interface
type clusterPartitionsEnv struct {
	partitions *cluster.Partitions
}

func (e *clusterPartitionsEnv) List() []PartitionEnv {
	clusterPartitions := e.partitions.Partitions()
	partitions := make([]PartitionEnv, len(clusterPartitions))
	for i, partition := range clusterPartitions {
		partitions[i] = &clusterPartitionEnv{
			clusterServiceEnv: &clusterServiceEnv{
				service: partition.Service,
			},
		}
	}
	return partitions
}

func (e *clusterPartitionsEnv) Get(name string) PartitionEnv {
	return &clusterPartitionEnv{
		clusterServiceEnv: &clusterServiceEnv{
			service: e.partitions.Partition(name).Service,
		},
	}
}

func (e *clusterPartitionsEnv) Connect() (*client.PartitionGroup, error) {
	return e.partitions.Connect()
}
