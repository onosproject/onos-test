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
	corev1 "k8s.io/api/core/v1"
)

// DatabaseSetup is an interface for setting up Raft partitions
type DatabaseSetup interface {
	Setup
	concurrentSetup
	Image(image string) DatabaseSetup
	PullPolicy(pullPolicy corev1.PullPolicy) DatabaseSetup
	Partitions(partitions int) DatabaseSetup
	Nodes(nodes int) DatabaseSetup
}

var _ DatabaseSetup = &databaseSetup{}

// databaseSetup is an implementation of the DatabaseSetup interface
type databaseSetup struct {
	*testSetup
	image      string
	pullPolicy corev1.PullPolicy
	partitions int
	nodes      int
}

func (s *databaseSetup) Image(image string) DatabaseSetup {
	s.image = image
	return s
}

func (s *databaseSetup) PullPolicy(pullPolicy corev1.PullPolicy) DatabaseSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *databaseSetup) Partitions(partitions int) DatabaseSetup {
	s.partitions = partitions
	return s
}

func (s *databaseSetup) Nodes(nodes int) DatabaseSetup {
	s.nodes = nodes
	return s
}

func (s *databaseSetup) create() error {
	return nil
}

func (s *databaseSetup) waitForStart() error {
	return nil
}
