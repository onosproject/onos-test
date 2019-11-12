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
	"fmt"
	"github.com/onosproject/onos-test/pkg/util/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"
)

// TopoType topology type
type TopoType int

const (
	// Single node topology
	Single TopoType = iota
	// Linear linear topology type
	Linear
	// Custom topology type
	Custom
)

func (d TopoType) String() string {
	return [...]string{"linear", "single"}[d]
}

func newNetwork(name string, client *client) *Network {
	return &Network{
		Node: newNode(name, 0, networkImage, client),
	}
}

// Network is an implementation of the Network interface
type Network struct {
	// TODO: Network should be a Service with a Node per device
	*Node
	topoType TopoType
	topo     string
	devices  int
}

// Devices returns a list of devices in the network
func (s *Network) Devices() []*Node {
	services, err := s.kubeClient.CoreV1().Services(s.namespace).List(metav1.ListOptions{
		LabelSelector: "type=network,network=" + s.name,
	})
	if err != nil {
		panic(err)
	}

	devices := make([]*Node, len(services.Items))
	for i, service := range services.Items {
		devices[i] = newNode(service.Name, 50001+i, "", s.client)
	}
	return devices
}

// SetSingle sets the network topology to single
func (s *Network) SetSingle() *Network {
	s.topoType = Single
	return s
}

// SetLinear sets the network to a linear topology
func (s *Network) SetLinear(devices int) *Network {
	s.topoType = Linear
	s.devices = devices
	return s
}

// SetTopo sets the network topology
func (s *Network) SetTopo(topo string, devices int) *Network {
	s.topoType = Custom
	s.topo = topo
	s.devices = devices
	return s
}

// Add adds the network to the cluster
func (s *Network) Add() error {
	step := logging.NewStep(s.namespace, fmt.Sprintf("Add network %s", s.Name()))
	step.Start()
	step.Logf("Creating %s Pod", s.Name())
	if err := s.createPod(); err != nil {
		step.Fail(err)
		return err
	}
	step.Logf("Creating %s Service", s.Name())
	if err := s.createService(); err != nil {
		step.Fail(err)
		return err
	}
	step.Logf("Waiting for %s to become ready", s.Name())
	if err := s.awaitReady(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// createPod creates a stratum Network pod
func (s *Network) createPod() error {
	var topoSpec string
	switch s.topoType {
	case Single:
		topoSpec = "single"
	case Linear:
		topoSpec = fmt.Sprintf("linear,%d", s.devices)
	case Custom:
		topoSpec = s.topo
	}

	var isPrivileged = true
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
			Labels: map[string]string{
				"type":    string(networkType),
				"Network": s.name,
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
func (s *Network) awaitReady() error {
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
func (s *Network) getNumDevices() int {
	switch s.topoType {
	case Single:
		return 1
	default:
		return s.devices
	}
}

// createService creates a Network service
func (s *Network) createService() error {
	var port int32 = 50001
	for i := 0; i < s.getNumDevices(); i++ {
		serviceName := fmt.Sprintf("%s-%d", s.name, i)
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: s.namespace,
				Labels: map[string]string{
					"type":    string(networkType),
					"network": s.name,
					"device":  serviceName,
				},
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"type":    string(networkType),
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

// Remove removes the network from the cluster
func (s *Network) Remove() error {
	step := logging.NewStep(s.namespace, fmt.Sprintf("Remove network %s", s.Name()))
	pods, err := s.kubeClient.CoreV1().Pods(s.namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("type=network,network=%s", s.name),
	})
	if err != nil {
		return err
	} else if len(pods.Items) == 0 {
		return fmt.Errorf("no resources matching '%s' found", s.name)
	}

	total := len(pods.Items)

	step.Logf("Deleting %s Pod", s.Name())
	if e := s.deletePod(); e != nil {
		err = e
	}
	if e := s.deleteService(); e != nil {
		err = e
	}

	for total > 0 {
		time.Sleep(50 * time.Millisecond)
		pods, err = s.kubeClient.CoreV1().Pods(s.namespace).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf("type=network,network=%s", s.name),
		})
		if err != nil {
			return err
		}

		total = len(pods.Items)

	}
	return err
}

// deletePod deletes a network Pod by name
func (s *Network) deletePod() error {
	return s.kubeClient.CoreV1().Pods(s.namespace).Delete(s.name, &metav1.DeleteOptions{})
}

// deleteService deletes all network Service by name
func (s *Network) deleteService() error {
	label := "type=network,network=" + s.name
	serviceList, _ := s.kubeClient.CoreV1().Services(s.namespace).List(metav1.ListOptions{
		LabelSelector: label,
	})

	for _, svc := range serviceList.Items {
		err := s.kubeClient.CoreV1().Services(s.namespace).Delete(svc.Name, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
