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

package env

import (
	"github.com/onosproject/onos-test/pkg/new/onit/setup"
	corev1 "k8s.io/api/core/v1"
)

// NetworkEnv provides the environment for a network node
type NetworkEnv interface {
	ServiceEnv
}

var _ NetworkEnv = &networkEnv{}

// networkEnv is an implementation of the NetworkEnv interface
type networkEnv struct {
	*serviceEnv
}

var _ setup.NetworkSetup = &networkSetup{}

// networkSetup is an implementation of the NetworkSetup interface
type networkSetup struct {
	*testEnv
	name       string
	image      string
	pullPolicy corev1.PullPolicy
}

func (s *networkSetup) Name(name string) setup.NetworkSetup {
	s.name = name
	return s
}

func (s *networkSetup) Image(image string) setup.NetworkSetup {
	s.image = image
	return s
}

func (s *networkSetup) PullPolicy(pullPolicy corev1.PullPolicy) setup.NetworkSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *networkSetup) Setup() error {
	return nil
}
