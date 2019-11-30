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

// PartitionsSetup is an interface for setting up Raft partitions
type PartitionsSetup interface {
	// Raft configures the partitions to use the Raft consensus protocol
	Raft() RaftPartitionsSetup

	// NOPaxos configures the partitions to use the NOPaxos consensus protocol
	NOPaxos() NOPaxosPartitionsSetup
}

var _ PartitionsSetup = &clusterPartitionsSetup{}

// clusterPartitionsSetup is an implementation of the PartitionsSetup interface
type clusterPartitionsSetup struct {
	group *cluster.Partitions
}

func (s *clusterPartitionsSetup) Raft() RaftPartitionsSetup {
	return &clusterRaftPartitionsSetup{
		raft: s.group.Raft(),
	}
}

func (s *clusterPartitionsSetup) NOPaxos() NOPaxosPartitionsSetup {
	return &clusterNOPaxosPartitionsSetup{
		nopaxos: s.group.NOPaxos(),
	}
}

func (s *clusterPartitionsSetup) setup() error {
	return s.group.Setup()
}
