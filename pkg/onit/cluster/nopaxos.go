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
	atomix "github.com/atomix/go-client/pkg/client"
	"github.com/atomix/kubernetes-controller/pkg/apis/cloud/v1beta1"
	"github.com/onosproject/onos-test/pkg/util/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	defaultReplicaImage   = "atomix/nopaxos-replica:latest"
	defaultSequencerImage = "atomix/nopaxos-proxy:latest"
)

func newNOPaxosPartitions(partitions *Partitions) *NOPaxosPartitions {
	return &NOPaxosPartitions{
		Partitions:          partitions,
		partitions:          1,
		replicas:            1,
		sequencerImage:      defaultSequencerImage,
		sequencerPullPolicy: corev1.PullIfNotPresent,
		replicaImage:        defaultReplicaImage,
		replicaPullPolicy:   corev1.PullIfNotPresent,
	}
}

// NOPaxosPartitions provides methods for adding and modifying NOPaxos partitions
type NOPaxosPartitions struct {
	*Partitions
	partitions          int
	replicas            int
	sequencerImage      string
	sequencerPullPolicy corev1.PullPolicy
	replicaImage        string
	replicaPullPolicy   corev1.PullPolicy
}

// NumPartitions returns the number of partitions
func (s *NOPaxosPartitions) NumPartitions() int {
	return GetArg(s.group, "partitions").Int(1)
}

// SetPartitions sets the number of partitions in the group
func (s *NOPaxosPartitions) SetPartitions(partitions int) {
	s.partitions = partitions
}

// Replicas returns the number of replicas in each partition
func (s *NOPaxosPartitions) Replicas() int {
	return GetArg(s.group, "replicas").Int(1)
}

// SetReplicas sets the number of nodes in each partition
func (s *NOPaxosPartitions) SetReplicas(replicas int) {
	s.replicas = replicas
}

// SequencerImage returns the image for the partition group
func (s *NOPaxosPartitions) SequencerImage() string {
	return GetArg(s.group, "sequencer", "image").String(s.sequencerImage)
}

// SetSequencerImage sets the image for the partition group sequencer
func (s *NOPaxosPartitions) SetSequencerImage(image string) {
	s.sequencerImage = image
}

// SequencerPullPolicy returns the image pull policy for the partition group
func (s *NOPaxosPartitions) SequencerPullPolicy() corev1.PullPolicy {
	return corev1.PullPolicy(GetArg(s.group, "sequencer", "pullPolicy").String(string(s.sequencerPullPolicy)))
}

// SetSequencerPullPolicy sets the image pull policy for the partition group
func (s *NOPaxosPartitions) SetSequencerPullPolicy(pullPolicy corev1.PullPolicy) {
	s.sequencerPullPolicy = pullPolicy
}

// ReplicaImage returns the image for the partition group
func (s *NOPaxosPartitions) ReplicaImage() string {
	return GetArg(s.group, "replica", "image").String(s.replicaImage)
}

// SetReplicaImage sets the image for the partition group replicas
func (s *NOPaxosPartitions) SetReplicaImage(image string) {
	s.replicaImage = image
}

// ReplicaPullPolicy returns the image pull policy for the partition group
func (s *NOPaxosPartitions) ReplicaPullPolicy() corev1.PullPolicy {
	return corev1.PullPolicy(GetArg(s.group, "replica", "pullPolicy").String(string(s.replicaPullPolicy)))
}

// SetReplicaPullPolicy sets the image pull policy for the partition group
func (s *NOPaxosPartitions) SetReplicaPullPolicy(pullPolicy corev1.PullPolicy) {
	s.replicaPullPolicy = pullPolicy
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
	database := &v1beta1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name(),
			Namespace: s.namespace,
		},
		Spec: v1beta1.DatabaseSpec{
			Clusters:   int32(s.NumPartitions()),
			Partitions: 1,
			Template: v1beta1.ClusterTemplateSpec{
				Spec: v1beta1.ClusterSpec{
					Proxy: &v1beta1.Proxy{
						Image:           s.SequencerImage(),
						ImagePullPolicy: s.SequencerPullPolicy(),
					},
					Backend: v1beta1.Backend{
						Replicas:        int32(s.Replicas()),
						Image:           s.ReplicaImage(),
						ImagePullPolicy: s.ReplicaPullPolicy(),
					},
				},
			},
		},
	}
	_, err := s.atomixClient.CloudV1beta1().Databases(s.namespace).Create(database)
	return err
}

// AwaitReady waits for partitions to complete startup
func (s *NOPaxosPartitions) AwaitReady() error {
	for {
		set, err := s.atomixClient.CloudV1beta1().Databases(s.namespace).Get(s.group, metav1.GetOptions{})
		if err != nil {
			return err
		} else if set.Status.ReadyPartitions == set.Spec.Partitions {
			return nil
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
