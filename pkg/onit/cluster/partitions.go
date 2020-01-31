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

package cluster

import (
	"context"
	"errors"
	atomix "github.com/atomix/go-client/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func newPartitions(cluster *Cluster, group string) *Partitions {
	return &Partitions{
		client:  cluster.client,
		cluster: cluster,
		group:   group,
	}
}

// Partitions provides methods for adding and modifying partitions
type Partitions struct {
	*client
	cluster *Cluster
	group   string
	raft    *RaftPartitions
	nopaxos *NOPaxosPartitions
}

// Name returns the partition group name
func (s *Partitions) Name() string {
	return s.group
}

// Raft returns new RaftPartitions
func (s *Partitions) Raft() *RaftPartitions {
	s.raft = newRaftPartitions(s)
	return s.raft
}

// NOPaxos returns new NOPaxosPartitions
func (s *Partitions) NOPaxos() *NOPaxosPartitions {
	s.nopaxos = newNOPaxosPartitions(s)
	return s.nopaxos
}

// Partition gets a partition by name
func (s *Partitions) Partition(name string) *Partition {
	return newPartition(s.cluster, name)
}

// getLabels returns the labels for the partition group
func (s *Partitions) getLabels() map[string]string {
	labels := getLabels(partitionType)
	labels[groupLabel] = s.group
	return labels
}

// Partitions lists the partitions in the group
func (s *Partitions) Partitions() []*Partition {
	labelSelector := metav1.LabelSelector{MatchLabels: s.getLabels()}
	list, err := s.atomixClient.CloudV1beta1().Partitions(s.namespace).List(metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String()})
	if err != nil {
		panic(err)
	}

	partitions := make([]*Partition, 0, len(list.Items))
	for _, partition := range list.Items {
		partitions = append(partitions, s.Partition(partition.Name))
	}
	return partitions
}

// Connect connects to the partition group
func (s *Partitions) Connect() (*atomix.Database, error) {
	client, err := atomix.NewClient("atomix-controller:5679", atomix.WithNamespace(s.namespace))
	if err != nil {
		return nil, err
	}
	return client.GetDatabase(context.Background(), s.group)
}

// Setup sets up a partition set
func (s *Partitions) Setup() error {
	if s.raft != nil {
		return s.raft.Setup()
	} else if s.nopaxos != nil {
		return s.nopaxos.Setup()
	}
	return errors.New("unknown partition type")
}

// AwaitReady waits for partitions to complete startup
func (s *Partitions) AwaitReady() error {
	if s.raft != nil {
		return s.raft.AwaitReady()
	} else if s.nopaxos != nil {
		return s.nopaxos.AwaitReady()
	}
	return errors.New("unknown partition type")
}
