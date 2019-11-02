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

// AtomixSetup is an interface for setting up the Atomix controller
type AtomixSetup interface {
	Setup
	sequentialSetup
	Image(image string) AtomixSetup
}

var _ AtomixSetup = &atomixSetup{}

// atomixSetup is an implementation of the AtomixSetup interface
type atomixSetup struct {
	*testSetup
	image      string
	pullPolicy corev1.PullPolicy
}

func (s *atomixSetup) Image(image string) AtomixSetup {
	s.image = image
	return s
}

func (s *atomixSetup) PullPolicy(pullPolicy corev1.PullPolicy) AtomixSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *atomixSetup) setup() error {
	return nil
}
