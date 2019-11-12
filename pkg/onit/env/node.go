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

import "github.com/onosproject/onos-test/pkg/onit/cluster"

// Node provides the environment for a single node
type Node interface {
	// Name returns the name of the node
	Name() string

	// Kill kills the node
	Kill() error

	// KillOrDie kills the node or panics if an error occurs
	KillOrDie()
}

// clusterNode is an implementation of the Node interface
type clusterNode struct {
	node *cluster.Node
}

func (e *clusterNode) Name() string {
	return e.node.Name()
}

func (e *clusterNode) Kill() error {
	return e.node.Delete()
}

func (e *clusterNode) KillOrDie() {
	if err := e.Kill(); err != nil {
		panic(err)
	}
}
