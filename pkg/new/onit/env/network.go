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
	"fmt"
	"github.com/onosproject/onos-test/pkg/new/onit/setup"
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

func newNetworkSetup(name string, testEnv *testEnv) setup.NetworkSetup {
	setup := &networkSetup{
		serviceSetup: &serviceSetup{
			testEnv: testEnv,
		},
		name: name,
	}
	setup.serviceSetup.setup = setup
	return setup
}

// Network provides the environment for a network node
type Network interface {
	Service

	// Devices returns a list of devices in the network
	Devices() []Node
}

var _ Network = &network{}

// network is an implementation of the Network interface
type network struct {
	*service
}

func (e *network) Devices() []Node {
	services, err := e.kubeClient.CoreV1().Services(e.namespace).List(metav1.ListOptions{
		LabelSelector: "type=network,network=" + e.name,
	})
	if err != nil {
		panic(err)
	}

	devices := make([]Node, len(services.Items))
	for i, service := range services.Items {
		devices[i] = &node{
			testEnv: e.testEnv,
			name:    service.Name,
		}
	}
	return devices
}

func (e *network) Remove() {
	if err := e.teardownNetwork(); err != nil {
		panic(err)
	}
}

// teardownNetwork tears down a network by name
func (e *network) teardownNetwork() error {
	pods, err := e.kubeClient.CoreV1().Pods(e.namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("type=network,network=%s", e.name),
	})
	if err != nil {
		return err
	} else if len(pods.Items) == 0 {
		return fmt.Errorf("no resources matching '%s' found", e.name)
	}
	total := len(pods.Items)
	if e := e.deletePod(); e != nil {
		err = e
	}
	if e := e.deleteService(); e != nil {
		err = e
	}

	for total > 0 {
		time.Sleep(50 * time.Millisecond)
		pods, err = e.kubeClient.CoreV1().Pods(e.namespace).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf("type=network,network=%s", e.name),
		})
		if err != nil {
			return err
		}

		total = len(pods.Items)

	}
	return err
}

// deletePod deletes a network Pod by name
func (e *network) deletePod() error {
	return e.kubeClient.CoreV1().Pods(e.namespace).Delete(e.name, &metav1.DeleteOptions{})
}

// deleteService deletes all network Service by name
func (e *network) deleteService() error {
	label := "type=network,network=" + e.name
	serviceList, _ := e.kubeClient.CoreV1().Services(e.namespace).List(metav1.ListOptions{
		LabelSelector: label,
	})

	for _, svc := range serviceList.Items {
		err := e.kubeClient.CoreV1().Services(e.namespace).Delete(svc.Name, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

var _ setup.NetworkSetup = &networkSetup{}

// networkSetup is an implementation of the NetworkSetup interface
type networkSetup struct {
	*serviceSetup
	name     string
	topoType TopoType
	nodes    int
}

func (s *networkSetup) Single() setup.NetworkSetup {
	s.topoType = Single
	return s
}

func (s *networkSetup) Linear(nodes int) setup.NetworkSetup {
	s.topoType = Single
	s.nodes = nodes
	return s
}

func (s *networkSetup) Using() setup.ServiceSetup {
	return s
}

func (s *networkSetup) Setup() error {
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
func (s *networkSetup) createPod() error {
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
func (s *networkSetup) awaitReady() error {
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
func (s *networkSetup) getNumDevices() int {
	switch s.topoType {
	case Single:
		return 1
	case Linear:
		return s.nodes
	}
	return 0
}

// createService creates a network service
func (s *networkSetup) createService() error {
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
