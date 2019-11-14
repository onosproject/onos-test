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

// CLISetup is an interface for setting up CLI nodes
type CLISetup interface {
	// Nodes sets the number of clusterTopo nodes to deploy
	Nodes(nodes int) CLISetup

	// Image sets the onos-cli image to deploy
	Image(image string) CLISetup

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) CLISetup
}

var _ CLISetup = &clusterCLISetup{}

// clusterCLISetup is an implementation of the CLI interface
type clusterCLISetup struct {
	cli *cluster.CLI
}

func (s *clusterCLISetup) Nodes(nodes int) CLISetup {
	s.cli.SetReplicas(nodes)
	return s
}

func (s *clusterCLISetup) Image(image string) CLISetup {
	s.cli.SetImage(image)
	return s
}

func (s *clusterCLISetup) PullPolicy(pullPolicy corev1.PullPolicy) CLISetup {
	s.cli.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterCLISetup) setup() error {
	return s.cli.Setup()
}
