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

// CLI is an interface for setting up CLI nodes
type CLI interface {
	// Image sets the onos-cli image to deploy
	Image(image string) CLI

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) CLI
}

var _ CLI = &clusterCLI{}

// clusterCLI is an implementation of the CLI interface
type clusterCLI struct {
	cli *cluster.CLI
}

func (s *clusterCLI) Image(image string) CLI {
	s.cli.SetImage(image)
	return s
}

func (s *clusterCLI) PullPolicy(pullPolicy corev1.PullPolicy) CLI {
	s.cli.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterCLI) create() error {
	return s.cli.Create()
}

// waitForStart waits for the onos-topo pods to complete startup
func (s *clusterCLI) waitForStart() error {
	return s.cli.AwaitReady()
}
