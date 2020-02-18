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

// SetVolume sets containers volumes
func (c *Container) SetVolume(volume ...VolumeMount) {
	c.volumes = volume
}

// Volume returns the container volume
func (c *Container) Volume() []VolumeMount {
	return c.volumes
}

// Name returns the container name
func (c *Container) Name() string {
	return GetArg(c.name, "service").String(c.name)
}

// Image returns the image for the container
func (c *Container) Image() string {
	return GetArg(c.name, "image").String(c.image)
}

// SetPullPolicy sets the image pull policy for the container
func (c *Container) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	c.pullPolicy = pullPolicy
}

// PullPolicy returns the image pull policy for the container
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

// Args returns the container arguments
func (c *Container) Args() []string {
	return c.args
}

// SetArgs sets the container arguments
func (c *Container) SetArgs(args ...string) {
	c.args = args
}

// SetCommand sets the container command
func (c *Container) SetCommand(command ...string) {
	c.command = command
}

// Command returns the container command
func (c *Container) Command() []string {
	return c.command
}

// Env returns the container environment variables
func (c *Container) Env() map[string]string {
	return c.env
}

// SetEnv sets the container environment variables
func (c *Container) SetEnv(env map[string]string) {
	c.env = env
}

// AddEnv adds an environment variable
func (c *Container) AddEnv(name, value string) {
	c.env[name] = value
}

// VolumeMount volume mount data structure
type VolumeMount struct {
	name string
	path string
}

// Container provides methods for adding sidecar containers to a deployment
type Container struct {
	name       string
	image      string
	pullPolicy corev1.PullPolicy
	command    []string
	env        map[string]string
	args       []string
	volumes    []VolumeMount
}
