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

import "time"

func newNodes(labels map[string]string, client *client) *Nodes {
	return &Nodes{
		client: client,
		labels: labels,
	}
}

// Nodes is a collection of nodes
type Nodes struct {
	*client
	labels map[string]string
}

// Get gets a node by name
func (n *Nodes) Get(name string) *Node {
	return newNode(name, "", n.client)
}

// List returns a list of nodes in the service
func (n *Nodes) List() []*Node {
	names := n.listPods(n.labels)
	nodes := make([]*Node, len(names))
	for i, name := range names {
		nodes[i] = n.Get(name)
	}
	return nodes
}

// AwaitReady waits for the nodes to become ready
func (n *Nodes) AwaitReady() error {
	for {
		ready, err := n.isReady()
		if err != nil {
			return err
		} else if ready {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// isReady returns a bool indicating whether all nodes are ready
func (n *Nodes) isReady() (bool, error) {
	nodes := n.List()
	for _, node := range nodes {
		if ready, err := node.isReady(); err != nil || !ready {
			return ready, err
		}
	}
	return true, nil
}
