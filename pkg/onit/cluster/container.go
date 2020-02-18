// Copyright 2020-present Open Networking Foundation.
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
)

func newContainer(cluster *Cluster) *Container {
	return &Container{
		pullPolicy: corev1.PullIfNotPresent,
	}
}

// Name returns the deployment name
func (c *Container) Name() string {
	return GetArg(c.name, "service").String(c.name)
}

// Image returns the image for the service
func (c *Container) Image() string {
	return GetArg(c.name, "image").String(c.image)
}

// SetPullPolicy sets the image pull policy for the service
func (c *Container) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	c.pullPolicy = pullPolicy
}

// PullPolicy returns the image pull policy for the service
func (c *Container) PullPolicy() corev1.PullPolicy {
	return corev1.PullPolicy(GetArg(c.name, "pullPolicy").String(string(c.pullPolicy)))
}

// SetName sets the container name
func (c *Container) SetName(name string) {
	c.name = name
}

// SetImage sets the container image
func (c *Container) SetImage(image string) {
	c.image = image
}

// Args returns the service arguments
func (c *Container) Args() []string {
	return c.args
}

// SetArgs sets the service arguments
func (c *Container) SetArgs(args ...string) {
	c.args = args
}

// SetCommand sets the service command
func (c *Container) SetCommand(command ...string) {
	c.command = command
}

// Command returns the service command
func (c *Container) Command() []string {
	return c.command
}

// Container provides methods for adding containers to a deployment
type Container struct {
	name       string
	image      string
	pullPolicy corev1.PullPolicy
	command    []string
	env        map[string]string
	args       []string
}
