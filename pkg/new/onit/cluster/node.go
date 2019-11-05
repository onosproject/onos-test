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

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newNode(name string, client *client) *Node {
	return &Node{
		client: client,
		name:   name,
	}
}

// Node provides the environment for a single node
type Node struct {
	*client
	name       string
	image      string
	pullPolicy corev1.PullPolicy
}

// Name returns the node name
func (n *Node) Name() string {
	return n.name
}

// SetName sets the node name
func (n *Node) SetName(name string) {
	n.name = name
}

// Image returns the image configured for the node
func (n *Node) Image() string {
	return n.image
}

// SetImage sets the node image
func (n *Node) SetImage(image string) {
	n.image = image
}

// PullPolicy returns the image pull policy configured for the node
func (n *Node) PullPolicy() corev1.PullPolicy {
	return n.pullPolicy
}

// SetPullPolicy sets the image pull policy for the node
func (n *Node) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	n.pullPolicy = pullPolicy
}

// Delete deletes the node
func (n *Node) Delete() error {
	return n.kubeClient.CoreV1().Pods(n.namespace).Delete(n.name, &metav1.DeleteOptions{})
}
