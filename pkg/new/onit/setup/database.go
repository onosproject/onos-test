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

package setup

// Database is an interface for setting up Raft partitions
type Database interface {
	ServiceType
	concurrentSetup

	// Partitions sets the number of partitions to deploy
	Partitions(partitions int) Database

	// Nodes sets the number of nodes per partition
	Nodes(nodes int) Database
}

var _ Database = &database{}

// database is an implementation of the Database interface
type database struct {
	*serviceType
	partitions int
	nodes      int
}

func (s *database) Partitions(partitions int) Database {
	s.partitions = partitions
	return s
}

func (s *database) Nodes(nodes int) Database {
	s.nodes = nodes
	return s
}

func (s *database) create() error {
	return nil
}

func (s *database) waitForStart() error {
	return nil
}
