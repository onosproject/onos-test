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
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
	corev1 "k8s.io/api/core/v1"
)

// Database is an interface for setting up Raft partitions
type Database interface {
	// Partitions sets the number of partitions to deploy
	Partitions(partitions int) Database

	// Nodes sets the number of nodes per partition
	Nodes(nodes int) Database

	// Image sets the Raft image to deploy
	Image(image string) Database

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) Database
}

var _ Database = &clusterDatabase{}

// clusterDatabase is an implementation of the Database interface
type clusterDatabase struct {
	group *cluster.Partitions
}

func (s *clusterDatabase) Partitions(partitions int) Database {
	s.group.SetPartitions(partitions)
	return s
}

func (s *clusterDatabase) Nodes(nodes int) Database {
	s.group.SetNodes(nodes)
	return s
}

func (s *clusterDatabase) Image(image string) Database {
	s.group.SetImage(image)
	return s
}

func (s *clusterDatabase) PullPolicy(pullPolicy corev1.PullPolicy) Database {
	s.group.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterDatabase) create() error {
	return s.group.Create()
}

func (s *clusterDatabase) waitForStart() error {
	return s.group.AwaitReady()
}
