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

package deploy

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"
)

// TopoType topology type
type TopoType int

const (
	// Linear linear topology type
	Linear TopoType = iota
	// Single node topology
	Single
)

func (d TopoType) String() string {
	return [...]string{"linear", "single"}[d]
}

// newNetworkDeploy returns a new Network deployer
func newNetworkDeploy(name string, deployment *testDeployment) Network {
	network := &network{
		serviceType: &serviceType{
			service: &service{
				testDeployment: deployment,
			},
		},
		name: name,
	}
	network.serviceType.deploy = network
	return network
}

// Network is an interface for deploying up a network
type Network interface {
	Deploy
	ServiceType

	// Single creates a single node topology
	Single() Network

	// Linear creates a linear topology with the given number of devices
	Linear(devices int) Network
}

var _ Network = &network{}

// network is an implementation of the Network interface
type network struct {
	*serviceType
	name     string
	topoType TopoType
	nodes    int
}

func (s *network) Single() Network {
	s.topoType = Single
	return s
}

func (s *network) Linear(nodes int) Network {
	s.topoType = Single
	s.nodes = nodes
	return s
}

func (s *network) Using() Service {
	return s
}

func (s *network) Setup() error {
	if err := s.createPod(); err != nil {
		return err
	}
	if err := s.createService(); err != nil {
		return err
	}
	if err := s.awaitReady(); err != nil {
		return err
	}
	return nil
}

// createPod creates a stratum network pod
func (s *network) createPod() error {
	var topoSpec string
	switch s.topoType {
	case Single:
		topoSpec = "single"
	case Linear:
		topoSpec = fmt.Sprintf("linear,%d", s.nodes)
	}

	var isPrivileged = true
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
			Labels: map[string]string{
				"type":    "network",
				"network": s.name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "stratum-simulator",
					Image:           s.image,
					ImagePullPolicy: s.pullPolicy,
					Stdin:           true,
					Args:            []string{"--topo", topoSpec},
					Ports: []corev1.ContainerPort{
						{
							Name:          "stratum",
							ContainerPort: 50001,
						},
					},
					ReadinessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(50001),
							},
						},
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(50001),
							},
						},
						InitialDelaySeconds: 15,
						PeriodSeconds:       20,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &isPrivileged,
					},
				},
			},
		},
	}

	_, err := s.kubeClient.CoreV1().Pods(s.namespace).Create(pod)
	return err
}

// awaitReady waits for the given simulator to complete startup
func (s *network) awaitReady() error {
	for {
		pod, err := s.kubeClient.CoreV1().Pods(s.namespace).Get(s.name, metav1.GetOptions{})
		if err != nil {
			return err
		} else if len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].Ready {
			return nil
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// getNumDevices returns the number of devices in the topology
func (s *network) getNumDevices() int {
	switch s.topoType {
	case Single:
		return 1
	case Linear:
		return s.nodes
	}
	return 0
}

// createService creates a network service
func (s *network) createService() error {
	var port int32 = 50001
	for i := 0; i < s.getNumDevices(); i++ {
		serviceName := fmt.Sprintf("%s-%d", s.name, i)
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: s.namespace,
				Labels: map[string]string{
					"type":    "network",
					"network": s.name,
					"device":  serviceName,
				},
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"type":    "network",
					"network": s.name,
				},
				Ports: []corev1.ServicePort{
					{
						Name: "stratum",
						Port: port,
					},
				},
			},
		}
		_, err := s.kubeClient.CoreV1().Services(s.namespace).Create(service)
		if err != nil {
			return err
		}
		port = port + 1
	}

	return nil
}
