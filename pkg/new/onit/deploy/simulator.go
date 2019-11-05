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
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	_, path, _, _     = runtime.Caller(0)
	deviceConfigsPath = filepath.Join(filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(path)))), "configs"), "device")
)

// newSimulatorDeploy returns a new Simulator deployer
func newSimulatorDeploy(name string, deployment *testDeployment) Simulator {
	simulator := &simulator{
		serviceType: &serviceType{
			service: &service{
				testDeployment: deployment,
			},
		},
		name: name,
	}
	simulator.serviceType.deploy = simulator
	return simulator
}

// Simulator is an interface for setting up a simulator
type Simulator interface {
	Deploy
	ServiceType
}

var _ Simulator = &simulator{}

// simulator is an implementation of the Simulator interface
type simulator struct {
	*serviceType
	name string
}

func (s *simulator) Using() Service {
	return s
}

func (s *simulator) Setup() error {
	if err := s.createConfigMap(); err != nil {
		return err
	}
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

// createConfigMap creates a simulator configuration
func (s *simulator) createConfigMap() error {
	file, err := os.Open(filepath.Join(deviceConfigsPath, "default.json"))
	if err != nil {
		return err
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
		},
		Data: map[string]string{
			"config.json": string(jsonBytes),
		},
	}
	_, err = s.kubeClient.CoreV1().ConfigMaps(s.namespace).Create(cm)
	return err
}

// createPod creates a simulator pod
func (s *simulator) createPod() error {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
			Labels: map[string]string{
				"type":      "simulator",
				"simulator": s.name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "onos-device-simulator",
					Image:           s.image,
					ImagePullPolicy: s.pullPolicy,
					Env: []corev1.EnvVar{
						{
							Name:  "GNMI_PORT",
							Value: "10161",
						},
						{
							Name:  "GNMI_INSECURE_PORT",
							Value: "11161",
						},
					},
					Ports: []corev1.ContainerPort{
						{
							Name:          "secure",
							ContainerPort: 10161,
						},
						{
							Name:          "insecure",
							ContainerPort: 11161,
						},
					},
					ReadinessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(11161),
							},
						},
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(11161),
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
func (s *simulator) createService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"type":      "simulator",
				"simulator": s.name,
			},
			Ports: []corev1.ServicePort{
				{
					Name: "secure",
					Port: 10161,
				},
				{
					Name: "insecure",
					Port: 11161,
				},
			},
		},
	}

	_, err := s.kubeClient.CoreV1().Services(s.namespace).Create(service)
	return err
}

// awaitReady waits for the given simulator to complete startup
func (s *simulator) awaitReady() error {
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
