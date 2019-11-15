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
	"github.com/onosproject/onos-test/pkg/onit/cluster"
)

// DeploymentEnv is a base interface for deployment environments
type DeploymentEnv interface {
	// Name is the name of the deployment
	Name() string

	// Nodes returns the deployment nodes
	Nodes() []NodeEnv

	// Node returns a specific node environment
	Node(name string) NodeEnv

	// AwaitReady waits for all nodes in the deployment to become ready
	AwaitReady() error
}

// clusterDeploymentEnv is an implementation of the Deployment interface
type clusterDeploymentEnv struct {
	deployment *cluster.Deployment
}

func (e *clusterDeploymentEnv) Name() string {
	return e.deployment.Name()
}

func (e *clusterDeploymentEnv) Nodes() []NodeEnv {
	clusterNodes, err := e.deployment.Nodes()
	if err != nil {
		panic(err)
	}
	nodes := make([]NodeEnv, len(clusterNodes))
	for i, node := range clusterNodes {
		nodes[i] = e.Node(node.Name())
	}
	return nodes
}

func (e *clusterDeploymentEnv) Node(name string) NodeEnv {
	node, err := e.deployment.Node(name)
	if err != nil {
		panic(err)
	}
	return &clusterNodeEnv{
		node,
	}
}

func (e *clusterDeploymentEnv) AwaitReady() error {
	return e.deployment.AwaitReady()
}
