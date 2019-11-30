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
	"fmt"
	atomix "github.com/atomix/atomix-go-client/pkg/client"
	"github.com/atomix/atomix-k8s-controller/pkg/apis/k8s/v1alpha1"
	"github.com/onosproject/onos-test/pkg/util/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func newNOPaxosPartitions(partitions *Partitions) *NOPaxosPartitions {
	return &NOPaxosPartitions{
		Partitions: partitions,
	}
}

// NOPaxosPartitions provides methods for adding and modifying NOPaxos partitions
type NOPaxosPartitions struct {
	*Partitions
	partitions     int
	replicas       int
	sequencerImage string
	replicaImage   string
	pullPolicy     corev1.PullPolicy
}

// SetPartitions sets the number of partitions in the group
func (s *NOPaxosPartitions) SetPartitions(partitions int) {
	s.partitions = partitions
}

// Replicas returns the number of replicas in each partition
func (s *NOPaxosPartitions) Replicas() int {
	return s.replicas
}

// SetReplicas sets the number of nodes in each partition
func (s *NOPaxosPartitions) SetReplicas(replicas int) {
	s.replicas = replicas
}

// SequencerImage returns the image for the partition group
func (s *NOPaxosPartitions) SequencerImage() string {
	return s.sequencerImage
}

// SetSequencerImage sets the image for the partition group sequencer
func (s *NOPaxosPartitions) SetSequencerImage(image string) {
	s.sequencerImage = image
}

// ReplicaImage returns the image for the partition group
func (s *NOPaxosPartitions) ReplicaImage() string {
	return s.replicaImage
}

// SetReplicaImage sets the image for the partition group replicas
func (s *NOPaxosPartitions) SetReplicaImage(image string) {
	s.replicaImage = image
}

// PullPolicy returns the image pull policy for the partition group
func (s *NOPaxosPartitions) PullPolicy() corev1.PullPolicy {
	return s.pullPolicy
}

// SetPullPolicy sets the image pull policy for the partition group
func (s *NOPaxosPartitions) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	s.pullPolicy = pullPolicy
}

// getLabels returns the labels for the partition group
func (s *NOPaxosPartitions) getLabels() map[string]string {
	labels := getLabels(partitionType)
	labels[groupLabel] = s.group
	return labels
}

// Connect connects to the partition group
func (s *NOPaxosPartitions) Connect() (*atomix.PartitionGroup, error) {
	client, err := atomix.NewClient(fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", s.namespace), atomix.WithNamespace(s.namespace))
	if err != nil {
		return nil, err
	}
	return client.GetGroup(context.Background(), s.group)
}

// Setup sets up a partition set
func (s *NOPaxosPartitions) Setup() error {
	step := logging.NewStep(s.namespace, "Setup NOPaxos partitions")
	step.Start()
	step.Log("Creating NOPaxos PartitionSet")
	if err := s.createPartitionSet(); err != nil {
		step.Fail(err)
		return err
	}
	step.Log("Waiting for NOPaxos partitions to become ready")
	if err := s.AwaitReady(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

func (s *NOPaxosPartitions) createPartitionSet() error {
	set := &v1alpha1.PartitionSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.group,
			Namespace: s.namespace,
		},
		Spec: v1alpha1.PartitionSetSpec{
			Partitions: s.partitions,
			Template: v1alpha1.PartitionTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: s.getLabels(),
				},
				Spec: v1alpha1.PartitionSpec{
					Size: int32(s.replicas),
					NOPaxos: &v1alpha1.NOPaxosProtocol{
						Sequencer: v1alpha1.NOPaxosSequencerSpec{
							Image:           s.sequencerImage,
							ImagePullPolicy: s.pullPolicy,
						},
						Protocol: v1alpha1.NOPaxosProtocolSpec{
							Image:           s.replicaImage,
							ImagePullPolicy: s.pullPolicy,
						},
					},
				},
			},
		},
	}
	_, err := s.atomixClient.K8sV1alpha1().PartitionSets(s.namespace).Create(set)
	return err
}

// AwaitReady waits for partitions to complete startup
func (s *NOPaxosPartitions) AwaitReady() error {
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
