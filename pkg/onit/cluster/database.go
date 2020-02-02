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

func newDatabase(cluster *Cluster, name string) *Database {
	return &Database{
		client:  cluster.client,
		cluster: cluster,
		name:    name,
	}
}

// Database provides methods for adding and modifying a database
type Database struct {
	*client
	cluster *Cluster
	name    string
	raft    *RaftDatabase
	nopaxos *NOPaxosDatabase
}

// Name returns the database name
func (s *Database) Name() string {
	return s.name
}

// Raft returns new RaftDatabase
func (s *Database) Raft() *RaftDatabase {
	s.raft = newRaftDatabase(s)
	return s.raft
}

// NOPaxos returns new NOPaxosDatabase
func (s *Database) NOPaxos() *NOPaxosDatabase {
	s.nopaxos = newNOPaxosDatabase(s)
	return s.nopaxos
}

// Partition gets a partition by name
func (s *Database) Partition(name string) *Partition {
	return newPartition(s.cluster, name)
}

// getLabels returns the labels for the database
func (s *Database) getLabels() map[string]string {
	labels := getLabels(partitionType)
	labels[groupLabel] = s.name
	return labels
}

// Database lists the partitions in the database
func (s *Database) Partitions() []*Partition {
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

// Connect connects to the database
func (s *Database) Connect() (*atomix.Database, error) {
	client, err := atomix.NewClient("atomix-controller:5679", atomix.WithNamespace(s.namespace))
	if err != nil {
		return nil, err
	}
	return client.GetDatabase(context.Background(), s.name)
}

// Setup sets up a database
func (s *Database) Setup() error {
	if s.raft != nil {
		return s.raft.Setup()
	} else if s.nopaxos != nil {
		return s.nopaxos.Setup()
	}
	return errors.New("unknown partition type")
}

// AwaitReady waits for database to complete startup
func (s *Database) AwaitReady() error {
	if s.raft != nil {
		return s.raft.AwaitReady()
	} else if s.nopaxos != nil {
		return s.nopaxos.AwaitReady()
	}
	return errors.New("unknown partition type")
}
