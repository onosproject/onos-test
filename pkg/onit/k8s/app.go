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
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GetApps returns a list of apps deployed in the cluster
func (c *ClusterController) GetApps() ([]string, error) {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": "onos", "type": "app"}}
	appList, err := c.kubeclient.AppsV1().Deployments(c.clusterID).List(metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String()})
	if err != nil {
		return nil, err
	}

	apps := make([]string, len(appList.Items))
	for i, app := range appList.Items {
		apps[i] = app.Name
	}
	return apps, nil
}

// setupApp creates an app
func (c *ClusterController) setupApp(name string, config *AppConfig) error {
	if err := c.createAppConfigMap(name, config); err != nil {
		return err
	}
	if err := c.createAppService(name); err != nil {
		return err
	}
	if err := c.createOnosAppDeployment(name, config.Image, config.PullPolicy); err != nil {
		return err
	}
	if err := c.awaitOnosAppDeploymentReady(name); err != nil {
		return err
	}
	return nil
}

// createAppConfigMap creates an app configuration
func (c *ClusterController) createAppConfigMap(name string, config *AppConfig) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.clusterID,
		},
		Data: map[string]string{},
	}
	_, err := c.kubeclient.CoreV1().ConfigMaps(c.clusterID).Create(cm)
	return err
}

// createOnosAppDeployment creates an app Deployment
func (c *ClusterController) createOnosAppDeployment(name string, image string, pullPolicy corev1.PullPolicy) error {
	nodes := int32(1)
	zero := int64(0)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.clusterID,
			Labels: map[string]string{
				"app":      "onos",
				"type":     "app",
				"resource": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nodes,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":      "onos",
					"type":     "app",
					"resource": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":      "onos",
						"type":     "app",
						"resource": name,
					},
					Annotations: map[string]string{
						"seccomp.security.alpha.kubernetes.io/pod": "unconfined",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            name,
							Image:           image,
							ImagePullPolicy: pullPolicy,
							Env: []corev1.EnvVar{
								{
									Name:  "ATOMIX_CONTROLLER",
									Value: fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", c.clusterID),
								},
								{
									Name:  "ATOMIX_APP",
									Value: "test",
								},
								{
									Name:  "ATOMIX_NAMESPACE",
									Value: c.clusterID,
								},
								{
									Name:  "ATOMIX_RAFT_GROUP",
									Value: "raft",
								},
							},
							Args: []string{
								"-caPath=/etc/app/certs/onf.cacrt",
								"-keyPath=/etc/app/certs/onos-config.key",
								"-certPath=/etc/app/certs/onos-config.crt",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "grpc",
									ContainerPort: 5150,
								},
								{
									Name:          "debug",
									ContainerPort: 40000,
									Protocol:      corev1.ProtocolTCP,
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
									Name:      "app",
									MountPath: "/etc/app/configs",
									ReadOnly:  true,
								},
								{
									Name:      "secret",
									MountPath: "/etc/app/certs",
									ReadOnly:  true,
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{
										"SYS_PTRACE",
									},
								},
							},
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser: &zero,
					},
					Volumes: []corev1.Volume{
						{
							Name: "app",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: name,
									},
								},
							},
						},
						{
							Name: "secret",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: c.clusterID,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Create(dep)
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
				"app":      "onos",
				"type":     "app",
				"resource": name,
			},
			Ports: []corev1.ServicePort{
				{
					Name: "grpc",
					Port: 5150,
				},
			},
		},
	}

	_, err := c.kubeclient.CoreV1().Services(c.clusterID).Create(service)
	return err
}

// awaitOnosAppDeploymentReady waits for the app pods to complete startup
func (c *ClusterController) awaitOnosAppDeploymentReady(name string) error {
	for {
		// Get the app deployment
		dep, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Return once the all replicas in the deployment are ready
		if int(dep.Status.ReadyReplicas) == c.config.ConfigNodes {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// teardownApp tears down a app by name
func (c *ClusterController) teardownApp(name string) error {
	var err error
	if e := c.deleteAppDeployment(name); e != nil {
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
func (c *ClusterController) deleteAppDeployment(name string) error {
	return c.kubeclient.AppsV1().Deployments(c.clusterID).Delete(name, &metav1.DeleteOptions{})
}

// deleteAppService deletes an app Service by name
func (c *ClusterController) deleteAppService(name string) error {
	return c.kubeclient.CoreV1().Services(c.clusterID).Delete(name, &metav1.DeleteOptions{})
}
