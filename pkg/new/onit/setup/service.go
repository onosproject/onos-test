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

import (
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
	corev1 "k8s.io/api/core/v1"
)

// ServiceType provides methods for setting up a service
type ServiceType interface {
	// Using returns the service setup
	Using() Service
}

// Service provides methods for setting up a service
type Service interface {
	// Image sets the simulator image to deploy
	Image(image string) Service

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) Service
}

// clusterServiceType is an implementation of the ServiceType interface
type clusterServiceType struct {
	*clusterService
}

func (s *clusterServiceType) Using() Service {
	return s.clusterService
}

var _ Service = &clusterService{}

// clusterService is an implementation of the Service interface
type clusterService struct {
	service *cluster.Service
}

func (s *clusterService) Image(image string) Service {
	s.service.SetImage(image)
	return s
}

func (s *clusterService) PullPolicy(pullPolicy corev1.PullPolicy) Service {
	s.service.SetPullPolicy(pullPolicy)
	return s
}
