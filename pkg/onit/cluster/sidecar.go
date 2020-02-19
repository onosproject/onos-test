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

func newSidecar(cluster *Cluster) *Sidecar {
	return &Sidecar{
		pullPolicy: corev1.PullIfNotPresent,
	}
}

// SetVolume sets sidecar containers volumes
func (s *Sidecar) SetVolume(volume ...corev1.VolumeMount) {
	s.volumes = volume
}

// Volume returns the sidecar container volume
func (s *Sidecar) Volume() []corev1.VolumeMount {
	return s.volumes
}

// Name returns the sidecar container name
func (s *Sidecar) Name() string {
	return GetArg(s.name, "service").String(s.name)
}

// Image returns the image for the sidecar container
func (s *Sidecar) Image() string {
	return GetArg(s.name, "image").String(s.image)
}

// SetPullPolicy sets the image pull policy for the container
func (s *Sidecar) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	s.pullPolicy = pullPolicy
}

// PullPolicy returns the image pull policy for the container
func (s *Sidecar) PullPolicy() corev1.PullPolicy {
	return corev1.PullPolicy(GetArg(s.name, "pullPolicy").String(string(s.pullPolicy)))
}

// SetName sets the container name
func (s *Sidecar) SetName(name string) {
	s.name = name
}

// SetImage sets the container image
func (s *Sidecar) SetImage(image string) {
	s.image = image
}

// Args returns the container arguments
func (s *Sidecar) Args() []string {
	return s.args
}

// SetArgs sets the sidecar container arguments
func (s *Sidecar) SetArgs(args ...string) {
	s.args = args
}

// SetCommand sets the sidecar container command
func (s *Sidecar) SetCommand(command ...string) {
	s.command = command
}

// Command returns the sidecar container command
func (s *Sidecar) Command() []string {
	return s.command
}

// Env returns the sidecar container environment variables
func (s *Sidecar) Env() map[string]string {
	return s.env
}

// SetEnv sets the sidecar container environment variables
func (s *Sidecar) SetEnv(env map[string]string) {
	s.env = env
}

// AddEnv adds an environment variable
func (s *Sidecar) AddEnv(name, value string) {
	s.env[name] = value
}

// Sidecar provides methods for adding sidecar containers to a deployment
type Sidecar struct {
	name       string
	image      string
	pullPolicy corev1.PullPolicy
	command    []string
	env        map[string]string
	args       []string
	volumes    []corev1.VolumeMount
}
