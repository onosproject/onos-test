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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	defaultRaftImage = "atomix/raft-replica:latest"
)

func newRaftPartitions(partitions *Partitions) *RaftPartitions {
	return &RaftPartitions{
		Partitions: partitions,
		partitions: 1,
		clusters:   1,
		replicas:   1,
		image:      defaultRaftImage,
		pullPolicy: corev1.PullIfNotPresent,
	}
}

// RaftPartitions provides methods for adding and modifying Raft partitions
type RaftPartitions struct {
	*Partitions
	partitions int
	clusters   int
	replicas   int
	image      string
	pullPolicy corev1.PullPolicy
}

// NumPartitions returns the number of partitions
func (s *RaftPartitions) NumPartitions() int {
	return GetArg(s.group, "partitions").Int(s.partitions)
}

// SetPartitions sets the number of partitions in the group
func (s *RaftPartitions) SetPartitions(partitions int) {
	s.partitions = partitions
}

// Clusters returns the number of clusters in each partition
func (s *RaftPartitions) Clusters() int {
	return GetArg(s.group, "clusters").Int(s.clusters)
}

// SetClusters sets the number of clusters in each partition
func (s *RaftPartitions) SetClusters(clusters int) {
	s.clusters = clusters
}

// Replicas returns the number of replicas in each partition
func (s *RaftPartitions) Replicas() int {
	return GetArg(s.group, "replicas").Int(s.replicas)
}

// SetReplicas sets the number of nodes in each partition
func (s *RaftPartitions) SetReplicas(replicas int) {
	s.replicas = replicas
}

// Image returns the image for the partition group
func (s *RaftPartitions) Image() string {
	return GetArg(s.group, "image").String(s.image)
}

// SetImage sets the image for the partition group
func (s *RaftPartitions) SetImage(image string) {
	s.image = image
}

// PullPolicy returns the image pull policy for the partition group
func (s *RaftPartitions) PullPolicy() corev1.PullPolicy {
	return corev1.PullPolicy(GetArg(s.group, "pullPolicy").String(string(s.pullPolicy)))
}

// SetPullPolicy sets the image pull policy for the partition group
func (s *RaftPartitions) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	s.pullPolicy = pullPolicy
}

// getLabels returns the labels for the partition group
func (s *RaftPartitions) getLabels() map[string]string {
	labels := getLabels(partitionType)
	labels[groupLabel] = s.group
	return labels
}

// Connect connects to the partition group
func (s *RaftPartitions) Connect() (*atomix.PartitionGroup, error) {
	client, err := atomix.NewClient(fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", s.namespace), atomix.WithNamespace(s.namespace))
	if err != nil {
		return nil, err
	}
	return client.GetGroup(context.Background(), s.group)
}

// Setup sets up a partition set
func (s *RaftPartitions) Setup() error {
	step := logging.NewStep(s.namespace, "Setup Raft partitions")
	step.Start()
	step.Log("Creating Raft Database")
	if err := s.createDatabase(); err != nil {
		step.Fail(err)
		return err
	}
	step.Log("Waiting for Raft partitions to become ready")
	if err := s.AwaitReady(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// createDatabase creates a Raft partition set from the configuration
func (s *RaftPartitions) createDatabase() error {
	var volumeClaim *corev1.PersistentVolumeClaim
	storageClass := GetArg(s.group, "storage", "class").String("")
	if storageClass != "" {
		storageSize := GetArg(s.group, "storage", "size").String("1G")
		volumeClaim = &corev1.PersistentVolumeClaim{
			Spec: corev1.PersistentVolumeClaimSpec{
				StorageClassName: &storageClass,
				AccessModes: []corev1.PersistentVolumeAccessMode{
					"ReadWriteOnce",
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						"storage": resource.MustParse(storageSize),
					},
				},
			},
		}
	}

	database := &v1beta1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name(),
			Namespace: s.namespace,
		},
		Spec: v1beta1.DatabaseSpec{
			Clusters: int32(s.Clusters()),
			Template: v1beta1.ClusterTemplateSpec{
				Spec: v1beta1.ClusterSpec{
					Partitions: int32(s.NumPartitions()),
					Backend: v1beta1.Backend{
						Replicas:        int32(s.Replicas()),
						Image:           s.Image(),
						ImagePullPolicy: s.PullPolicy(),
						VolumeClaim:     volumeClaim,
					},
				},
			},
		},
	}
	_, err := s.atomixClient.CloudV1beta1().Databases(s.namespace).Create(database)
	return err
}

// AwaitReady waits for partitions to complete startup
func (s *RaftPartitions) AwaitReady() error {
	for {
		database, err := s.atomixClient.CloudV1beta1().Databases(s.namespace).Get(s.group, metav1.GetOptions{})
		if err != nil {
			return err
		} else if database.Status.ReadyClusters == database.Spec.Clusters {
			return nil
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
