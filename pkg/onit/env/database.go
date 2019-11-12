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
	"github.com/onosproject/onos-test/pkg/onit/cluster"
)

// Database provides the database environment
type Database interface {
	// Partitions returns all database partitions
	Partitions(group string) Partitions
}

var _ Database = &clusterDatabase{}

// clusterDatabase is an implementation of the Database interface
type clusterDatabase struct {
	database *cluster.Database
}

func (e *clusterDatabase) Partitions(group string) Partitions {
	return &clusterPartitions{
		partitions: e.database.Partitions(group),
	}
}
