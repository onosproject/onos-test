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
	corev1 "k8s.io/api/core/v1"
)

// RaftDatabaseSetup is an interface for setting up Raft partitions
type RaftDatabaseSetup interface {
	// SetPartitions sets the number of partitions to deploy
	SetPartitions(partitions int) RaftDatabaseSetup

	// SetClusters sets the number of clusters in the database
	SetClusters(clusters int) RaftDatabaseSetup

	// SetReplicas sets the number of replicas per partition
	SetReplicas(replicas int) RaftDatabaseSetup

	// SetImage sets the Raft image to deploy
	SetImage(image string) RaftDatabaseSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) RaftDatabaseSetup
}

var _ RaftDatabaseSetup = &clusterRaftDatabaseSetup{}

// clusterRaftDatabaseSetup is an implementation of the RaftDatabaseSetup interface
type clusterRaftDatabaseSetup struct {
	raft *cluster.RaftDatabase
}

func (s *clusterRaftDatabaseSetup) SetPartitions(partitions int) RaftDatabaseSetup {
	s.raft.SetPartitions(partitions)
	return s
}

func (s *clusterRaftDatabaseSetup) SetClusters(clusters int) RaftDatabaseSetup {
	s.raft.SetClusters(clusters)
	return s
}

func (s *clusterRaftDatabaseSetup) SetReplicas(replicas int) RaftDatabaseSetup {
	s.raft.SetReplicas(replicas)
	return s
}

func (s *clusterRaftDatabaseSetup) SetImage(image string) RaftDatabaseSetup {
	s.raft.SetImage(image)
	return s
}

func (s *clusterRaftDatabaseSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) RaftDatabaseSetup {
	s.raft.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterRaftDatabaseSetup) setup() error {
	return s.raft.Setup()
}
