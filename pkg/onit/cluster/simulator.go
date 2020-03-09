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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"k8s.io/apimachinery/pkg/watch"

	"github.com/onosproject/onos-test/pkg/util/logging"
	"github.com/onosproject/onos-topo/api/device"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	simulatorType             = "simulator"
	simulatorLabel            = "simulator"
	simulatorImage            = "onosproject/device-simulator:latest"
	simulatorService          = "device-simulator"
	simulatorDeviceType       = "Devicesim"
	simulatorDeviceVersion    = "1.0.0"
	simulatorSecurePortName   = "secure"
	simulatorSecurePort       = 10161
	simulatorInsecurePortName = "insecure"
	simulatorInsecurePort     = 11161
	gnmiPortEnv               = "GNMI_PORT"
	gnmiInsecurePortEnv       = "GNMI_INSECURE_PORT"
)

func newSimulator(cluster *Cluster, name string) *Simulator {
	node := newNode(cluster)
	node.SetName(name)
	node.SetImage(simulatorImage)
	node.SetPort(simulatorInsecurePort)
	return &Simulator{
		Node:          node,
		add:           true,
		deviceType:    simulatorDeviceType,
		deviceVersion: simulatorDeviceVersion,
	}
}

// Simulator provides methods for adding and modifying simulators
type Simulator struct {
	*Node
	add           bool
	deviceType    string
	deviceVersion string
	deviceTimeout *time.Duration
}

// AddDevice returns whether to add the device to the topo service
func (s *Simulator) AddDevice() bool {
	return s.add
}

// SetAddDevice sets whether to add the device to the topo service
func (s *Simulator) SetAddDevice(add bool) {
	s.add = add
}

// DeviceType returns the device type
func (s *Simulator) DeviceType() string {
	return s.deviceType
}

// SetDeviceType sets the device type
func (s *Simulator) SetDeviceType(deviceType string) {
	s.deviceType = deviceType
}

// DeviceVersion returns the device version
func (s *Simulator) DeviceVersion() string {
	return s.deviceVersion
}

// SetDeviceVersion sets the device version
func (s *Simulator) SetDeviceVersion(version string) {
	s.deviceVersion = version
}

// DeviceTimeout returns the device timeout
func (s *Simulator) DeviceTimeout() *time.Duration {
	return s.deviceTimeout
}

// SetDeviceTimeout sets the device timeout
func (s *Simulator) SetDeviceTimeout(timeout time.Duration) {
	s.deviceTimeout = &timeout
}

// Setup adds the simulator to the cluster
func (s *Simulator) Setup() error {
	step := logging.NewStep(s.namespace, fmt.Sprintf("Add simulator %s", s.Name()))
	step.Start()
	step.Logf("Creating %s ConfigMap", s.Name())
	if err := s.createConfigMap(); err != nil {
		step.Fail(err)
		return err
	}
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
		step.Logf("Adding %s to onos-topo", s.Name())
		if err := s.addDevice(); err != nil {
			step.Fail(err)
			return err
		}
	}
	step.Complete()
	return nil
}

// getLabels gets the simulator labels
func (s *Simulator) getLabels() map[string]string {
	labels := getLabels(simulatorType)
	labels[simulatorLabel] = s.name
	return labels
}

// createConfigMap creates a simulator configuration
func (s *Simulator) createConfigMap() error {
	configMapsAPI := s.kubeClient.CoreV1().ConfigMaps("kube-test")
	configMap, err := configMapsAPI.Get("device-simulator", metav1.GetOptions{})
	if err != nil {
		return err
	}

	config := NewSimulatorConfig([]byte(configMap.Data["config.yaml"]))
	err = config.Parse()
	if err != nil {
		return err
	}
	configData := config.GetConfig()
	var simConfig []byte
	simConfig, err = json.Marshal(configData.Configuration)
	if err != nil {
		return err
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
			Labels:    s.getLabels(),
		},
		Data: map[string]string{
			"config.json": string(simConfig),
		},
	}
	_, err = s.kubeClient.CoreV1().ConfigMaps(s.namespace).Create(cm)
	return err

}

// createPod creates a simulator pod
func (s *Simulator) createPod() error {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
			Labels:    s.getLabels(),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            simulatorService,
					Image:           s.image,
					ImagePullPolicy: s.pullPolicy,
					Env: []corev1.EnvVar{
						{
							Name:  gnmiPortEnv,
							Value: fmt.Sprintf("%d", simulatorSecurePort),
						},
						{
							Name:  gnmiInsecurePortEnv,
							Value: fmt.Sprintf("%d", s.Port()),
						},
					},
					Ports: []corev1.ContainerPort{
						{
							Name:          simulatorSecurePortName,
							ContainerPort: simulatorSecurePort,
						},
						{
							Name:          simulatorInsecurePortName,
							ContainerPort: int32(s.Port()),
						},
					},
					ReadinessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(s.Port()),
							},
						},
						PeriodSeconds:    1,
						FailureThreshold: 30,
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(s.Port()),
							},
						},
						InitialDelaySeconds: 15,
						PeriodSeconds:       20,
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "config",
							MountPath: "/etc/simulator/configs",
							ReadOnly:  true,
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "config",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: s.name,
							},
						},
					},
				},
			},
		},
	}

	_, err := s.kubeClient.CoreV1().Pods(s.namespace).Create(pod)
	return err
}

// createService creates a simulator service
func (s *Simulator) createService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
			Labels:    s.getLabels(),
		},
		Spec: corev1.ServiceSpec{
			Selector: s.getLabels(),
			Ports: []corev1.ServicePort{
				{
					Name: simulatorSecurePortName,
					Port: simulatorSecurePort,
				},
				{
					Name: simulatorInsecurePortName,
					Port: int32(s.Port()),
				},
			},
		},
	}

	_, err := s.kubeClient.CoreV1().Services(s.namespace).Create(service)
	return err
}

// awaitReady waits for the given simulator to complete startup
func (s *Simulator) awaitReady() error {
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

// connectTopo connects to the topo service
func (s *Simulator) connectTopo() (*grpc.ClientConn, error) {
	tlsConfig, err := s.Credentials()
	if err != nil {
		return nil, err
	}
	return grpc.Dial(s.cluster.Topo().Address(s.cluster.Topo().Ports()[0].Name), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}

// addDevice adds the device to the topo service
func (s *Simulator) addDevice() error {
	tlsConfig, err := s.Credentials()
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(s.cluster.Topo().Address(s.cluster.Topo().Ports()[0].Name), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := device.NewDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), topoTimeout)
	defer cancel()
	_, err = client.Add(ctx, &device.AddRequest{
		Device: &device.Device{
			ID:      device.ID(s.Name()),
			Address: s.Address(),
			Type:    device.Type(s.DeviceType()),
			Version: s.DeviceVersion(),
			Timeout: s.deviceTimeout,
			TLS: device.TlsConfig{
				Plain: true,
			},
		},
	})
	return err
}

// AwaitDevicePredicate waits for the given device predicate
func (s *Simulator) AwaitDevicePredicate(predicate func(*device.Device) bool, timeout time.Duration) error {
	conn, err := s.connectTopo()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := device.NewDeviceServiceClient(conn)

	// Set a timer within which the device must reach the connected/available state
	errCh := make(chan error)
	timer := time.NewTimer(5 * time.Second)

	// Open a stream to listen for events from the device service
	stream, err := client.List(context.Background(), &device.ListRequest{
		Subscribe: true,
	})
	if err != nil {
		return err
	}

	// Start a goroutine to listen for device events from the topo service
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				errCh <- err
				close(errCh)
				return
			}

			if predicate(response.Device) {
				timer.Stop()
				close(errCh)
				return
			}
		}
	}()

	select {
	// If the timer fires, return a timeout error
	case _, ok := <-timer.C:
		if !ok {
			return errors.New("device predicate timed out")
		}
		return nil
	// If an error is received on the error channel, return the error. Otherwise, return nil
	case err := <-errCh:
		return err
	}
}

// TearDown removes the simulator from the cluster
func (s *Simulator) TearDown() error {
	var err error

	if e := s.removeDevice(); e != nil {
		err = e
	}
	if e := s.deletePod(); e != nil {
		err = e
	}
	if e := s.deleteService(); e != nil {
		err = e
	}
	if e := s.deleteConfigMap(); e != nil {
		err = e
	}
	if e := s.waitForDeletePod(); e != nil {
		err = e
	}
	return err
}

// removeDevice removes the device from the topo service
func (s *Simulator) removeDevice() error {
	if !s.add {
		return nil
	}

	tlsConfig, err := s.Credentials()
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(s.cluster.Topo().Address(s.cluster.Topo().Ports()[0].Name), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := device.NewDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), topoTimeout)
	response, err := client.Get(ctx, &device.GetRequest{
		ID: device.ID(s.Name()),
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

// deleteConfigMap deletes a simulator ConfigMap by name
func (s *Simulator) deleteConfigMap() error {
	return s.kubeClient.CoreV1().ConfigMaps(s.namespace).Delete(s.name, &metav1.DeleteOptions{})
}

// deletePod deletes a simulator Pod by name
func (s *Simulator) deletePod() error {
	return s.kubeClient.CoreV1().Pods(s.namespace).Delete(s.name, &metav1.DeleteOptions{})
}

func (s *Simulator) waitForDeletePod() error {
	w, e := s.kubeClient.CoreV1().Pods(s.client.namespace).Watch(metav1.ListOptions{})
	if e != nil {
		return e
	}
	for event := range w.ResultChan() {
		switch event.Type {
		case watch.Deleted:
			w.Stop()
		}
	}
	return nil
}

// deleteService deletes a simulator Service by name
func (s *Simulator) deleteService() error {
	return s.kubeClient.CoreV1().Services(s.namespace).Delete(s.name, &metav1.DeleteOptions{})
}
