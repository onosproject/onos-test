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

// DatabaseSetup is an interface for setting up Raft partitions
type DatabaseSetup interface {
	// Partitions sets the number of partitions to deploy
	Partitions(partitions int) DatabaseSetup

	// Nodes sets the number of nodes per partition
	Nodes(nodes int) DatabaseSetup

	// Image sets the Raft image to deploy
	Image(image string) DatabaseSetup

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) DatabaseSetup
}

var _ DatabaseSetup = &clusterDatabaseSetup{}

// clusterDatabaseSetup is an implementation of the Database interface
type clusterDatabaseSetup struct {
	group *cluster.Partitions
}

func (s *clusterDatabaseSetup) Partitions(partitions int) DatabaseSetup {
	s.group.SetPartitions(partitions)
	return s
}

func (s *clusterDatabaseSetup) Nodes(nodes int) DatabaseSetup {
	s.group.SetNodes(nodes)
	return s
}

func (s *clusterDatabaseSetup) Image(image string) DatabaseSetup {
	s.group.SetImage(image)
	return s
}

func (s *clusterDatabaseSetup) PullPolicy(pullPolicy corev1.PullPolicy) DatabaseSetup {
	s.group.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterDatabaseSetup) setup() error {
	return s.group.Setup()
}
