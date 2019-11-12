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
func (s *Nodes) Get(name string) *Node {
	return newNode(name, "", s.client)
}

// List returns a list of nodes in the service
func (s *Nodes) List() []*Node {
	names := s.listPods(s.labels)
	nodes := make([]*Node, len(names))
	for i, name := range names {
		nodes[i] = s.Get(name)
	}
	return nodes
}
