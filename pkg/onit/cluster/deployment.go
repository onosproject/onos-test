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
	"errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const apiPort = "api"

func newDeployment(name string, labels map[string]string, image string, client *client) *Deployment {
	return &Deployment{
		client:     client,
		name:       name,
		labels:     labels,
		image:      image,
		pullPolicy: corev1.PullIfNotPresent,
	}
}

// Deployment is a collection of nodes
type Deployment struct {
	*client
	name       string
	labels     map[string]string
	image      string
	pullPolicy corev1.PullPolicy
}

// Name returns the deployment name
func (d *Deployment) Name() string {
	return GetArg(d.name, "service").String(d.name)
}

// Image returns the image for the service
func (d *Deployment) Image() string {
	return GetArg(d.name, "image").String(d.image)
}

// SetImage sets the image for the service
func (d *Deployment) SetImage(image string) {
	d.image = image
}

// PullPolicy returns the image pull policy for the service
func (d *Deployment) PullPolicy() corev1.PullPolicy {
	return corev1.PullPolicy(GetArg(d.name, "pullPolicy").String(string(d.pullPolicy)))
}

// SetPullPolicy sets the image pull policy for the service
func (d *Deployment) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	d.pullPolicy = pullPolicy
}

// Node gets a node by name
func (d *Deployment) Node(name string) (*Node, error) {
	pod, err := d.kubeClient.CoreV1().Pods(d.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			if port.Name == apiPort {
				return newNode(name, int(port.ContainerPort), container.Image, d.client), nil
			}
		}
	}
	return newNode(name, 0, "", d.client), nil
}

// Nodes returns a list of nodes in the service
func (d *Deployment) Nodes() ([]*Node, error) {
	names := d.listPods(d.labels)
	nodes := make([]*Node, len(names))
	for i, name := range names {
		node, err := d.Node(name)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	return nodes, nil
}

// AwaitReady waits for the nodes to become ready
func (d *Deployment) AwaitReady() error {
	for {
		ready, err := d.isReady()
		if err != nil {
			return err
		} else if ready {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// isReady returns a bool indicating whether all nodes are ready
func (d *Deployment) isReady() (bool, error) {
	nodes, err := d.Nodes()
	if err != nil {
		return false, err
	}
	for _, node := range nodes {
		if ready, err := node.isReady(); err != nil || !ready {
			return ready, err
		}
	}
	return true, nil
}

// Execute executes the given command on one of the service nodes
func (d *Deployment) Execute(command ...string) ([]string, int, error) {
	nodes, err := d.Nodes()
	if err != nil {
		return nil, 0, err
	}
	if len(nodes) == 0 {
		return nil, 0, errors.New("no nodes found")
	}
	return nodes[0].Execute(command...)
}
