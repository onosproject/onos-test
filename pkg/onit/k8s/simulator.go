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

package k8s

import (
	"encoding/json"
	"time"

	"github.com/onosproject/onos-test/pkg/onit/console"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GetSimulators returns a list of simulators deployed in the cluster
func (c *ClusterController) GetSimulators() ([]string, error) {
	pods, err := c.Kubeclient.CoreV1().Pods(c.ClusterID).List(metav1.ListOptions{
		LabelSelector: "type=simulator",
	})

	if err != nil {
		return nil, err
	}

	simulators := make([]string, len(pods.Items))
	for i, pod := range pods.Items {
		simulators[i] = pod.Name
	}
	return simulators, nil
}

// setupSimulator creates a simulator required for the test
func (c *ClusterController) setupSimulator(name string, config *SimulatorConfig) error {
	if err := c.createSimulatorConfigMap(name, config); err != nil {
		return err
	}
	if err := c.createSimulatorPod(name); err != nil {
		return err
	}
	if err := c.createSimulatorService(name); err != nil {
		return err
	}
	if err := c.awaitSimulatorReady(name); err != nil {
		return err
	}
	return nil
}

// createSimulatorConfigMap creates a simulator configuration
func (c *ClusterController) createSimulatorConfigMap(name string, config *SimulatorConfig) error {
	configObj, err := config.load()
	if err != nil {
		return err
	}
	configJSON, err := json.Marshal(configObj)
	if err != nil {
		return err
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.ClusterID,
		},
		Data: map[string]string{
			"config.json": string(configJSON),
		},
	}
	_, err = c.Kubeclient.CoreV1().ConfigMaps(c.ClusterID).Create(cm)
	return err
}

// createSimulatorPod creates a simulator pod
func (c *ClusterController) createSimulatorPod(name string) error {

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.ClusterID,
			Labels: map[string]string{
				"type":      "simulator",
				"simulator": name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "onos-device-simulator",
					Image:           c.imageName("onosproject/device-simulator", c.Config.ImageTags["simulator"]),
					ImagePullPolicy: corev1.PullIfNotPresent,
					Ports: []corev1.ContainerPort{
						{
							Name:          "gnmi",
							ContainerPort: 10161,
						},
					},
					ReadinessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(10161),
							},
						},
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(10161),
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
								Name: name,
							},
						},
					},
				},
			},
		},
	}

	_, err := c.Kubeclient.CoreV1().Pods(c.ClusterID).Create(pod)
	return err
}

// createSimulatorService creates a simulator service
func (c *ClusterController) createSimulatorService(name string) error {

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.ClusterID,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"simulator": name,
			},
			Ports: []corev1.ServicePort{
				{
					Name: "gnmi",
					Port: 10161,
				},
			},
		},
	}

	_, err := c.Kubeclient.CoreV1().Services(c.ClusterID).Create(service)
	return err
}

// awaitSimulatorReady waits for the given simulator to complete startup
func (c *ClusterController) awaitSimulatorReady(name string) error {
	for {
		pod, err := c.Kubeclient.CoreV1().Pods(c.ClusterID).Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		} else if len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].Ready {
			return nil
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// AddSimulator adds a device simulator with the given configuration
func (c *ClusterController) AddSimulator(name string, config *SimulatorConfig) console.ErrorStatus {
	c.Status.Start("Setting up simulator")
	if err := c.setupSimulator(name, config); err != nil {
		return c.Status.Fail(err)
	}
	c.Status.Start("Reconfiguring onos-config nodes")
	if err := c.addSimulatorToConfig(name); err != nil {
		return c.Status.Fail(err)
	}
	return c.Status.Succeed()
}

// RemoveSimulator removes a device simulator with the given name
func (c *ClusterController) RemoveSimulator(name string) console.ErrorStatus {
	c.Status.Start("Tearing down simulator")
	if err := c.teardownSimulator(name); err != nil {
		c.Status.Fail(err)
	}
	c.Status.Start("Reconfiguring onos-config nodes")
	if err := c.removeSimulatorFromConfig(name); err != nil {
		return c.Status.Fail(err)
	}
	return c.Status.Succeed()
}

// teardownSimulator tears down a simulator by name
func (c *ClusterController) teardownSimulator(name string) error {
	var err error
	if e := c.deleteSimulatorPod(name); e != nil {
		err = e
	}
	if e := c.deleteSimulatorService(name); e != nil {
		err = e
	}
	if e := c.deleteSimulatorConfigMap(name); e != nil {
		err = e
	}
	return err
}

// deleteSimulatorConfigMap deletes a simulator ConfigMap by name
func (c *ClusterController) deleteSimulatorConfigMap(name string) error {
	return c.Kubeclient.CoreV1().ConfigMaps(c.ClusterID).Delete(name, &metav1.DeleteOptions{})
}

// deleteSimulatorPod deletes a simulator Pod by name
func (c *ClusterController) deleteSimulatorPod(name string) error {
	return c.Kubeclient.CoreV1().Pods(c.ClusterID).Delete(name, &metav1.DeleteOptions{})
}

// deleteSimulatorService deletes a simulator Service by name
func (c *ClusterController) deleteSimulatorService(name string) error {
	return c.Kubeclient.CoreV1().Services(c.ClusterID).Delete(name, &metav1.DeleteOptions{})
}
