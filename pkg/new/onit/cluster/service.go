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

import corev1 "k8s.io/api/core/v1"

func newService(name string, serviceType serviceType, client *client) *Service {
	return &Service{
		client:      client,
		name:        name,
		serviceType: serviceType,
	}
}

// Service is the base type for multi-node services
type Service struct {
	*client
	name        string
	replicas    int
	serviceType serviceType
	image       string
	pullPolicy  corev1.PullPolicy
}

// Name returns the name of the service
func (s *Service) Name() string {
	return s.name
}

// SetName sets the service name
func (s *Service) SetName(name string) {
	s.name = name
}

// Nodes returns the collection of nodes in the service
func (s *Service) Nodes() *Nodes {
	return newNodes(s.name, s.serviceType, s.client)
}

// Replicas returns the number of nodes in the service
func (s *Service) Replicas() int {
	return s.replicas
}

// SetReplicas sets the number of nodes in the service
func (s *Service) SetReplicas(replicas int) {
	s.replicas = replicas
}

// Image returns the image for the service
func (s *Service) Image() string {
	return s.image
}

// SetImage sets the image for the service
func (s *Service) SetImage(image string) {
	s.image = image
}

// PullPolicy returns the image pull policy for the service
func (s *Service) PullPolicy() corev1.PullPolicy {
	return s.pullPolicy
}

// SetPullPolicy sets the image pull policy for the service
func (s *Service) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	s.pullPolicy = pullPolicy
}
