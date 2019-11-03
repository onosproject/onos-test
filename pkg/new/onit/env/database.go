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

// Database provides the database environment
type Database interface {
	// Partitions returns all database partitions
	Partitions() []Partition

	// Partition returns the Partition for the given partition ID
	Partition(id int) Partition
}

var _ Database = &database{}

// database is an implementation of the Database interface
type database struct {
	*testEnv
}

func (e *database) Partitions() []Partition {
	panic("implement me")
}

func (e database) Partition(id int) Partition {
	return &partition{
		service: &service{
			testEnv: e.testEnv,
		},
		id: id,
	}
}

func (e *database) Nodes(partition int) []Node {
	panic("implement me")
}
