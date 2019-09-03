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
	"time"

	"k8s.io/apimachinery/pkg/labels"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// setupGUI sets up the GUI resources
func (c *ClusterController) setupGUI() error {
	if err := c.createGUIDeployment(); err != nil {
		return err
	}
	if err := c.createGUIService(); err != nil {
		return err
	}
	if err := c.awaitGUIDeploymentReady(); err != nil {
		return err
	}
	return nil
}

// createGUIDeployment creates an onos-gui deployment
func (c *ClusterController) createGUIDeployment() error {
	nodes := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-gui",
			Namespace: c.clusterID,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nodes,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":      "onos",
					"type":     "gui",
					"resource": "onos-gui",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":      "onos",
						"type":     "gui",
						"resource": "onos-gui",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "onos-gui",
							Image:           c.imageName("onosproject/onos-gui", c.config.ImageTags["gui"]),
							ImagePullPolicy: c.config.PullPolicy,
							Ports: []corev1.ContainerPort{
								{
									Name:          "grpc",
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Create(deployment)
	return err
}

// createGUIService creates an onos-gui service
func (c *ClusterController) createGUIService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-gui",
			Namespace: c.clusterID,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":      "onos",
				"type":     "gui",
				"resource": "onos-gui",
			},
			Ports: []corev1.ServicePort{
				{
					Name: "grpc",
					Port: 80,
				},
			},
		},
	}
	_, err := c.kubeclient.CoreV1().Services(c.clusterID).Create(service)
	return err
}

// awaitGUIDeploymentReady waits for the onos-config proxy pods to complete startup
func (c *ClusterController) awaitGUIDeploymentReady() error {
	for {
		// Get the onos-gui deployment
		dep, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Get("onos-gui", metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Return once the all replicas in the deployment are ready
		if int(dep.Status.ReadyReplicas) == 1 {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// GetOnosGuiNodes returns a list of all onos-gui nodes running in the cluster
func (c *ClusterController) GetOnosGuiNodes() ([]NodeInfo, error) {
	configLabelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": "onos", "type": "gui"}}

	pods, err := c.kubeclient.CoreV1().Pods(c.clusterID).List(metav1.ListOptions{
		LabelSelector: labels.Set(configLabelSelector.MatchLabels).String(),
	})
	if err != nil {
		return nil, err
	}

	onosGuiNodes := make([]NodeInfo, len(pods.Items))
	for i, pod := range pods.Items {
		var status NodeStatus
		if pod.Status.Phase == corev1.PodRunning {
			status = NodeRunning
		} else if pod.Status.Phase == corev1.PodFailed {
			status = NodeFailed
		}
		onosGuiNodes[i] = NodeInfo{
			ID:     pod.Name,
			Status: status,
			Type:   OnosGui,
		}
	}

	return onosGuiNodes, nil
}
