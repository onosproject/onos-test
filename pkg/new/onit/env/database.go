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

// DatabaseEnv provides the database environment
type DatabaseEnv interface {
	// Partitions returns all database partitions
	Partitions() []PartitionEnv

	// Partition returns the PartitionEnv for the given partition ID
	Partition(id int) PartitionEnv
}

var _ DatabaseEnv = &databaseEnv{}

// databaseEnv is an implementation of the DatabaseEnv interface
type databaseEnv struct {
	*testEnv
}

func (e *databaseEnv) Partitions() []PartitionEnv {
	panic("implement me")
}

func (e databaseEnv) Partition(id int) PartitionEnv {
	return &partitionEnv{
		serviceEnv: &serviceEnv{
			testEnv: e.testEnv,
		},
		id: id,
	}
}

func (e *databaseEnv) Nodes(partition int) []NodeEnv {
	panic("implement me")
}
