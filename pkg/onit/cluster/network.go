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
	"context"
	"errors"
	"fmt"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"
)

const (
	networkType          = "network"
	networkLabel         = "network"
	networkImage         = "opennetworking/mn-stratum:latest"
	networkDeviceType    = "Stratum"
	networkDeviceVersion = "1.0.0"
	stratumPortName      = "stratum"
	stratumPort          = 28000
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
		Node:          newNode(name, 0, networkImage, client),
		add:           true,
		deviceType:    networkDeviceType,
		deviceVersion: networkDeviceVersion,
	}
}

// Network is an implementation of the Network interface
type Network struct {
	// TODO: Network should be a Service with a Node per device
	*Node
	topoType      TopoType
	topo          string
	devices       int
	add           bool
	deviceType    string
	deviceVersion string
	deviceTimeout *time.Duration
}

// Devices returns a list of devices in the network
func (s *Network) Devices() ([]*Node, error) {
	services, err := s.kubeClient.CoreV1().Services(s.namespace).List(metav1.ListOptions{
		LabelSelector: "type=network,network=" + s.name,
	})
	if err != nil {
		return nil, err
	}

	devices := make([]*Node, len(services.Items))
	for i, service := range services.Items {
		devices[i] = newNode(service.Name, stratumPort, "", s.client)
	}
	return devices, nil
}

// AddDevices returns whether to add the devices to the topo service
func (s *Network) AddDevices() bool {
	return s.add
}

// SetAddDevices sets whether to add the devices to the topo service
func (s *Network) SetAddDevices(add bool) {
	s.add = add
}

// DeviceType returns the device type
func (s *Network) DeviceType() string {
	return s.deviceType
}

// SetDeviceType sets the device type
func (s *Network) SetDeviceType(deviceType string) {
	s.deviceType = deviceType
}

// DeviceVersion returns the device version
func (s *Network) DeviceVersion() string {
	return s.deviceVersion
}

// SetDeviceVersion sets the device version
func (s *Network) SetDeviceVersion(version string) {
	s.deviceVersion = version
}

// DeviceTimeout returns the device timeout
func (s *Network) DeviceTimeout() *time.Duration {
	return s.deviceTimeout
}

// SetDeviceTimeout sets the device timeout
func (s *Network) SetDeviceTimeout(timeout time.Duration) {
	s.deviceTimeout = &timeout
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

// Setup sets up the network
func (s *Network) Setup() error {
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
	if s.add {
		step.Logf("Adding %s devices to topo service", s.Name())
		if err := s.addDevices(); err != nil {
			step.Fail(err)
			return err
		}
	}
	step.Complete()
	return nil
}

// getLabels gets the network labels
func (s *Network) getLabels() map[string]string {
	labels := getLabels(networkType)
	labels[networkLabel] = s.name
	return labels
}

// getDeviceNames returns a list of device names
func (s *Network) getDeviceNames() []string {
	numDevices := s.getNumDevices()
	names := make([]string, numDevices)
	for i := 0; i < numDevices; i++ {
		names[i] = fmt.Sprintf("%s-%d", s.name, i)
	}
	return names
}

// getDevicePorts returns a map of device names and ports
func (s *Network) getDevicePorts() map[string]int32 {
	names := s.getDeviceNames()
	ports := make(map[string]int32)
	var port int32 = 50001
	for _, name := range names {
		ports[name] = port
		port++
	}
	return ports
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

	devices := s.getDevicePorts()
	ports := make([]corev1.ContainerPort, 0, len(devices))
	for device, port := range s.getDevicePorts() {
		ports = append(ports, corev1.ContainerPort{
			Name:          device,
			ContainerPort: port,
		})
	}

	var isPrivileged = true
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
			Labels:    s.getLabels(),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "stratum-simulator",
					Image:           s.image,
					ImagePullPolicy: s.pullPolicy,
					Stdin:           true,
					Args:            []string{"--topo", topoSpec},
					Ports:           ports,
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
	for device, port := range s.getDevicePorts() {
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      device,
				Namespace: s.namespace,
				Labels:    s.getLabels(),
			},
			Spec: corev1.ServiceSpec{
				Selector: s.getLabels(),
				Ports: []corev1.ServicePort{
					{
						Name:       stratumPortName,
						Port:       stratumPort,
						TargetPort: intstr.FromInt(int(port)),
					},
				},
			},
		}
		_, err := s.kubeClient.CoreV1().Services(s.namespace).Create(service)
		if err != nil {
			return err
		}
	}
	return nil
}

// addDevices adds the network's devices to the topo service
func (s *Network) addDevices() error {
	devices, err := s.Devices()
	if err != nil {
		return err
	}
	for _, device := range devices {
		if err := s.addDevice(device); err != nil {
			return err
		}
	}
	return nil
}

// addDevice adds the given device to the topo service
func (s *Network) addDevice(node *Node) error {
	if err := s.addDeviceByCLI(node); err == nil {
		return nil
	}
	return s.addDeviceByAPI(node)
}

func (s *Network) addDeviceByCLI(node *Node) error {
	// Determine whether any CLI nodes are deployed and use the CLI to add the device if possible
	cli := newCLI(s.client)
	nodes, err := cli.Nodes()
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return errors.New("onos-cli is not available")
	}

	timeout := s.DeviceTimeout()
	if timeout == nil {
		t := topoTimeout
		timeout = &t
	}
	_, _, err = nodes[0].Execute(fmt.Sprintf("onos topo add device %s --address %s --type %s --version %s --timeout %s --plain", node.Name(), node.Address(), s.DeviceType(), s.DeviceVersion(), timeout))
	return err
}

func (s *Network) addDeviceByAPI(node *Node) error {
	tlsConfig, err := s.Credentials()
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(topoAddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), topoTimeout)
	defer cancel()
	client := device.NewDeviceServiceClient(conn)
	_, err = client.Add(ctx, &device.AddRequest{
		Device: &device.Device{
			ID:      device.ID(node.Name()),
			Address: node.Address(),
			Type:    device.Type(s.deviceType),
			Version: s.deviceVersion,
			Timeout: s.deviceTimeout,
			TLS: device.TlsConfig{
				Plain: true,
			},
		},
	})
	return err
}

// TearDown removes the network from the cluster
func (s *Network) TearDown() error {
	step := logging.NewStep(s.namespace, fmt.Sprintf("Remove network %s", s.Name()))

	var err error
	step.Logf("Removing %s devices", s.Name())
	if e := s.removeDevices(); e != nil {
		err = e
	}
	step.Logf("Deleting %s Pod", s.Name())
	if e := s.deletePod(); e != nil {
		err = e
	}
	step.Logf("Deleting %s Service", s.Name())
	if e := s.deleteService(); e != nil {
		err = e
	}
	step.Complete()
	return err
}

// removeDevices removes the devices from the topo service
func (s *Network) removeDevices() error {
	devices, err := s.Devices()
	if err != nil {
		return err
	}
	for _, device := range devices {
		if e := s.removeDevice(device); e != nil {
			err = e
		}
	}
	return err
}

// removeDevice removes the given device from the topo service
func (s *Network) removeDevice(node *Node) error {
	if err := s.removeDeviceByCLI(node); err == nil {
		return nil
	}
	return s.removeDeviceByAPI(node)
}

func (s *Network) removeDeviceByCLI(node *Node) error {
	// Determine whether any CLI nodes are deployed and use the CLI to add the device if possible
	cli := newCLI(s.client)
	nodes, err := cli.Nodes()
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return errors.New("onos-cli is not available")
	}
	_, _, err = nodes[0].Execute(fmt.Sprintf("onos topo remove device %s", node.Name()))
	return err
}

func (s *Network) removeDeviceByAPI(node *Node) error {
	tlsConfig, err := s.Credentials()
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(topoAddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := device.NewDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), topoTimeout)
	response, err := client.Get(ctx, &device.GetRequest{
		ID: device.ID(node.Name()),
	})
	cancel()
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), topoTimeout)
	_, err = client.Remove(ctx, &device.RemoveRequest{
		Device: response.Device,
	})
	cancel()
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
