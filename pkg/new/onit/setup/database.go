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
	"github.com/atomix/atomix-api/proto/atomix/protocols/raft"
	"github.com/atomix/atomix-k8s-controller/pkg/apis/k8s/v1alpha1"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

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

// create creates a Raft partition set
func (s *database) create() error {
	if err := s.createPartitionSet(); err != nil {
		return err
	}
	return nil
}

// createPartitionSet creates a Raft partition set from the configuration
func (s *database) createPartitionSet() error {
	bytes, err := yaml.Marshal(&raft.RaftProtocol{})
	if err != nil {
		return err
	}

	set := &v1alpha1.PartitionSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "raft",
			Namespace: s.namespace,
		},
		Spec: v1alpha1.PartitionSetSpec{
			Partitions: s.partitions,
			Template: v1alpha1.PartitionTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"type":  "database",
						"group": "raft",
					},
				},
				Spec: v1alpha1.PartitionSpec{
					Size:            int32(s.nodes),
					Protocol:        "raft",
					Image:           s.image,
					ImagePullPolicy: s.pullPolicy,
					Config:          string(bytes),
				},
			},
		},
	}
	_, err = s.atomixClient.K8sV1alpha1().PartitionSets(s.namespace).Create(set)
	return err
}

// waitForStart waits for Raft partitions to complete startup
func (s *database) waitForStart() error {
	for {
		set, err := s.atomixClient.K8sV1alpha1().PartitionSets(s.namespace).Get("raft", metav1.GetOptions{})
		if err != nil {
			return err
		} else if int(set.Status.ReadyPartitions) == set.Spec.Partitions {
			return nil
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
