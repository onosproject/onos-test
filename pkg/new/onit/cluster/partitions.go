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
	"github.com/atomix/atomix-api/proto/atomix/protocols/raft"
	"github.com/atomix/atomix-k8s-controller/pkg/apis/k8s/v1alpha1"
	"github.com/ghodss/yaml"
	"github.com/onosproject/onos-test/pkg/new/util/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"time"
)

func newPartitions(group string, client *client) *Partitions {
	return &Partitions{
		client: client,
		group:  group,
	}
}

// Partitions provides methods for adding and modifying partitions
type Partitions struct {
	*client
	group      string
	partitions int
	nodes      int
	image      string
	pullPolicy corev1.PullPolicy
}

// SetPartitions sets the number of partitions in the group
func (s *Partitions) SetPartitions(partitions int) {
	s.partitions = partitions
}

// Nodes returns the number of nodes in each partition
func (s *Partitions) Nodes() int {
	return s.nodes
}

// SetNodes sets the number of nodes in each partition
func (s *Partitions) SetNodes(nodes int) {
	s.nodes = nodes
}

// Image returns the image for the partition group
func (s *Partitions) Image() string {
	return s.image
}

// SetImage sets the image for the partition group
func (s *Partitions) SetImage(image string) {
	s.image = image
}

// PullPolicy returns the image pull policy for the partition group
func (s *Partitions) PullPolicy() corev1.PullPolicy {
	return s.pullPolicy
}

// SetPullPolicy sets the image pull policy for the partition group
func (s *Partitions) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	s.pullPolicy = pullPolicy
}

// Partition gets a partition by name
func (s *Partitions) Partition(name string) *Partition {
	return newPartition(name, s.client)
}

// Partitions lists the partitions in the group
func (s *Partitions) Partitions() []*Partition {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{typeLabel: databaseType.name(), "group": s.group}}
	list, err := s.atomixClient.K8sV1alpha1().Partitions(s.namespace).List(metav1.ListOptions{
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

// Create creates a partition set
func (s *Partitions) Create() error {
	step := logging.NewStep(s.namespace, "Create Raft partitions")
	step.Start()
	step.Log("Creating Raft PartitionSet")
	if err := s.createPartitionSet(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// createPartitionSet creates a Raft partition set from the configuration
func (s *Partitions) createPartitionSet() error {
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
						"type":  databaseType.name(),
						"group": s.group,
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

// AwaitReady waits for partitions to complete startup
func (s *Partitions) AwaitReady() error {
	for {
		set, err := s.atomixClient.K8sV1alpha1().PartitionSets(s.namespace).Get(s.group, metav1.GetOptions{})
		if err != nil {
			return err
		} else if int(set.Status.ReadyPartitions) == set.Spec.Partitions {
			return nil
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
