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

import (
	"github.com/onosproject/onos-test/pkg/onit/cluster"
)

// DatabaseSetup is an interface for setting up a database
type DatabaseSetup interface {
	// Raft configures the partitions to use the Raft consensus protocol
	Raft() RaftDatabaseSetup

	// NOPaxos configures the partitions to use the NOPaxos consensus protocol
	NOPaxos() NOPaxosDatabaseSetup
}

var _ DatabaseSetup = &clusterDatabaseSetup{}

// clusterDatabaseSetup is an implementation of the DatabaseSetup interface
type clusterDatabaseSetup struct {
	group *cluster.Database
}

func (s *clusterDatabaseSetup) Raft() RaftDatabaseSetup {
	return &clusterRaftDatabaseSetup{
		raft: s.group.Raft(),
	}
}

func (s *clusterDatabaseSetup) NOPaxos() NOPaxosDatabaseSetup {
	return &clusterNOPaxosPartitionsSetup{
		nopaxos: s.group.NOPaxos(),
	}
}

func (s *clusterDatabaseSetup) setup() error {
	return s.group.Setup()
}
