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

package env

import (
	"github.com/onosproject/onos-test/pkg/new/onit/setup"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// Service is a base interface for service environments
type Service interface {
	// Name is the name of the service
	Name() string

	// Nodes returns the service nodes
	Nodes() []Node

	// Node returns a specific node environment
	Node(name string) Node

	// Remove removes the service
	Remove()
}

// ServiceSetup is a base interface for services that can be set up
type ServiceSetup interface {
	Service
	setup.Setup
}

var _ Service = &service{}

func newService(name string, serviceType string, testEnv *testEnv) *service {
	return &service{
		testEnv:     testEnv,
		name:        name,
		serviceType: serviceType,
	}
}

// service is an implementation of the Service interface
type service struct {
	*testEnv
	name        string
	serviceType string
}

func (e *service) Name() string {
	return e.name
}

func (e *service) Nodes() []Node {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"type": e.serviceType}}
	pods, err := e.kubeClient.CoreV1().Pods(e.namespace).List(metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String()})
	if err != nil {
		panic(err)
	}

	nodes := make([]Node, len(pods.Items))
	for i, pod := range pods.Items {
		nodes[i] = e.Node(pod.Name)
	}
	return nodes
}

func (e *service) Node(name string) Node {
	return &node{
		testEnv: e.testEnv,
		name:    name,
	}
}

func (e *service) Remove() {
	panic("implement me")
}

var _ setup.ServiceSetup = &serviceSetup{}

// serviceSetup is an implementation of the ServiceSetup interface
type serviceSetup struct {
	*testEnv
	setup      setup.Setup
	image      string
	pullPolicy corev1.PullPolicy
}

func (s *serviceSetup) Image(image string) setup.ServiceSetup {
	s.image = image
	return s
}

func (s *serviceSetup) PullPolicy(pullPolicy corev1.PullPolicy) setup.ServiceSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *serviceSetup) Setup() error {
	return s.setup.Setup()
}

func (s *serviceSetup) SetupOrDie() {
	if err := s.Setup(); err != nil {
		panic(err)
	}
}
