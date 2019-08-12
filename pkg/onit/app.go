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

package onit

import (
	"encoding/json"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GetSApps returns a list of apps deployed in the cluster
func (c *ClusterController) GetApps() ([]string, error) {
	pods, err := c.kubeclient.CoreV1().Pods(c.clusterID).List(metav1.ListOptions{
		LabelSelector: "type=app",
	})

	if err != nil {
		return nil, err
	}

	apps := make([]string, len(pods.Items))
	for i, pod := range pods.Items {
		apps[i] = pod.Name
	}
	return apps, nil
}

// setupApp creates an app
func (c *ClusterController) setupApp(name string, config *AppConfig) error {
	if err := c.createAppConfigMap(name, config); err != nil {
		return err
	}
	if err := c.createAppPod(name); err != nil {
		return err
	}
	if err := c.createAppService(name); err != nil {
		return err
	}
	if err := c.awaitAppReady(name); err != nil {
		return err
	}
	return nil
}

// createAppConfigMap creates an app configuration
func (c *ClusterController) createAppConfigMap(name string, config *AppConfig) error {
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
			Namespace: c.clusterID,
		},
		Data: map[string]string{
			"config.json": string(configJSON),
		},
	}
	_, err = c.kubeclient.CoreV1().ConfigMaps(c.clusterID).Create(cm)
	return err
}

// createAppPod creates an app pod
func (c *ClusterController) createAppPod(name string) error {

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.clusterID,
			Labels: map[string]string{
				"type": "app",
				"app":  name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "app",
					// TODO: pull in the image name from the command arguments
					Image:           c.imageName("onosproject/onos-ztp", "latest"),
					ImagePullPolicy: corev1.PullIfNotPresent,
					Ports: []corev1.ContainerPort{
						{
							Name:          "gnmi",
							ContainerPort: 5150,
						},
					},
					ReadinessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(5150),
							},
						},
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(5150),
							},
						},
						InitialDelaySeconds: 15,
						PeriodSeconds:       20,
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "config",
							MountPath: "/etc/app/configs",
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

	_, err := c.kubeclient.CoreV1().Pods(c.clusterID).Create(pod)
	return err
}

// createAppService creates an app service
func (c *ClusterController) createAppService(name string) error {

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.clusterID,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": name,
			},
			Ports: []corev1.ServicePort{
				{
					Name: "gnmi",
					Port: 5150,
				},
			},
		},
	}

	_, err := c.kubeclient.CoreV1().Services(c.clusterID).Create(service)
	return err
}

// awaitAppReady waits for the given app to complete startup
func (c *ClusterController) awaitAppReady(name string) error {
	for {
		pod, err := c.kubeclient.CoreV1().Pods(c.clusterID).Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		} else if len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].Ready {
			return nil
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// teardownApp tears down a app by name
func (c *ClusterController) teardownApp(name string) error {
	var err error
	if e := c.deleteAppPod(name); e != nil {
		err = e
	}
	if e := c.deleteAppService(name); e != nil {
		err = e
	}
	if e := c.deleteAppConfigMap(name); e != nil {
		err = e
	}
	return err
}

// deleteAppConfigMap deletes an app ConfigMap by name
func (c *ClusterController) deleteAppConfigMap(name string) error {
	return c.kubeclient.CoreV1().ConfigMaps(c.clusterID).Delete(name, &metav1.DeleteOptions{})
}

// deleteAppPod deletes an app Pod by name
func (c *ClusterController) deleteAppPod(name string) error {
	return c.kubeclient.CoreV1().Pods(c.clusterID).Delete(name, &metav1.DeleteOptions{})
}

// deleteAppService deletes an app Service by name
func (c *ClusterController) deleteAppService(name string) error {
	return c.kubeclient.CoreV1().Services(c.clusterID).Delete(name, &metav1.DeleteOptions{})
}
