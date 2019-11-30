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

// RaftPartitionsSetup is an interface for setting up Raft partitions
type RaftPartitionsSetup interface {
	// SetPartitions sets the number of partitions to deploy
	SetPartitions(partitions int) RaftPartitionsSetup

	// SetReplicasPerPartition sets the number of replicas per partition
	SetReplicasPerPartition(replicas int) RaftPartitionsSetup

	// SetImage sets the Raft image to deploy
	SetImage(image string) RaftPartitionsSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) RaftPartitionsSetup
}

var _ RaftPartitionsSetup = &clusterRaftPartitionsSetup{}

// clusterRaftPartitionsSetup is an implementation of the RaftPartitionsSetup interface
type clusterRaftPartitionsSetup struct {
	raft *cluster.RaftPartitions
}

func (s *clusterRaftPartitionsSetup) SetPartitions(partitions int) RaftPartitionsSetup {
	s.raft.SetPartitions(partitions)
	return s
}

func (s *clusterRaftPartitionsSetup) SetReplicasPerPartition(replicas int) RaftPartitionsSetup {
	s.raft.SetReplicas(replicas)
	return s
}

func (s *clusterRaftPartitionsSetup) SetImage(image string) RaftPartitionsSetup {
	s.raft.SetImage(image)
	return s
}

func (s *clusterRaftPartitionsSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) RaftPartitionsSetup {
	s.raft.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterRaftPartitionsSetup) setup() error {
	return s.raft.Setup()
}
