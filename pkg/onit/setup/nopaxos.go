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

// NOPaxosPartitionsSetup is an interface for setting up NOPaxos partitions
type NOPaxosPartitionsSetup interface {
	// SetPartitions sets the number of partitions to deploy
	SetPartitions(partitions int) NOPaxosPartitionsSetup

	// SetReplicasPerPartition sets the number of replicas per partition
	SetReplicasPerPartition(replicas int) NOPaxosPartitionsSetup

	// SetSequencerImage sets the sequencer image to deploy
	SetSequencerImage(image string) NOPaxosPartitionsSetup

	// SetSequencerPullPolicy sets the sequencer image pull policy
	SetSequencerPullPolicy(pullPolicy corev1.PullPolicy) NOPaxosPartitionsSetup

	// SetReplicaImage sets the replica image to deploy
	SetReplicaImage(image string) NOPaxosPartitionsSetup

	// SetReplicaPullPolicy sets the replica image pull policy
	SetReplicaPullPolicy(pullPolicy corev1.PullPolicy) NOPaxosPartitionsSetup
}

var _ NOPaxosPartitionsSetup = &clusterNOPaxosPartitionsSetup{}

// clusterNOPaxosPartitionsSetup is an implementation of the NOPaxosPartitionsSetup interface
type clusterNOPaxosPartitionsSetup struct {
	nopaxos *cluster.NOPaxosPartitions
}

func (s *clusterNOPaxosPartitionsSetup) SetPartitions(partitions int) NOPaxosPartitionsSetup {
	s.nopaxos.SetPartitions(partitions)
	return s
}

func (s *clusterNOPaxosPartitionsSetup) SetReplicasPerPartition(replicas int) NOPaxosPartitionsSetup {
	s.nopaxos.SetReplicas(replicas)
	return s
}

func (s *clusterNOPaxosPartitionsSetup) SetSequencerImage(image string) NOPaxosPartitionsSetup {
	s.nopaxos.SetSequencerImage(image)
	return s
}

func (s *clusterNOPaxosPartitionsSetup) SetSequencerPullPolicy(pullPolicy corev1.PullPolicy) NOPaxosPartitionsSetup {
	s.nopaxos.SetSequencerPullPolicy(pullPolicy)
	return s
}

func (s *clusterNOPaxosPartitionsSetup) SetReplicaImage(image string) NOPaxosPartitionsSetup {
	s.nopaxos.SetReplicaImage(image)
	return s
}

func (s *clusterNOPaxosPartitionsSetup) SetReplicaPullPolicy(pullPolicy corev1.PullPolicy) NOPaxosPartitionsSetup {
	s.nopaxos.SetReplicaPullPolicy(pullPolicy)
	return s
}

func (s *clusterNOPaxosPartitionsSetup) setup() error {
	return s.nopaxos.Setup()
}
