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

package setup

import corev1 "k8s.io/api/core/v1"

// ServiceType provides methods for setting up a service
type ServiceType interface {
	// Using returns the generic service setup
	Using() Service
}

// ServiceTypeSetup provides methods for setting up a service
type ServiceTypeSetup interface {
	// Using returns the generic service setup
	Using() ServiceSetup
}

// Service provides methods for setting up a service
type Service interface {
	// Image sets the simulator image to deploy
	Image(image string) Service

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) Service
}

// ServiceSetup provides methods for setting up a service
type ServiceSetup interface {
	Setup

	// Image sets the simulator image to deploy
	Image(image string) ServiceSetup

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) ServiceSetup
}

var _ ServiceType = &serviceType{}

// serviceType is an implementation of the ServiceType interface
type serviceType struct {
	*service
}

func (s *serviceType) Using() Service {
	return s.service
}

var _ Service = &service{}

// service is an implementation of the Service interface
type service struct {
	*testSetup
	image      string
	pullPolicy corev1.PullPolicy
}

func (s *service) Image(image string) Service {
	s.image = image
	return s
}

func (s *service) PullPolicy(pullPolicy corev1.PullPolicy) Service {
	s.pullPolicy = pullPolicy
	return s
}
