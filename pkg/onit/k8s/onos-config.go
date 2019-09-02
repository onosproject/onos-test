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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// setupOnosConfig sets up the onos-config Deployment
func (c *ClusterController) setupOnosConfig() error {
	if err := c.createOnosConfigConfigMap(); err != nil {
		return err
	}
	if err := c.createOnosConfigService(); err != nil {
		return err
	}
	if err := c.createOnosConfigDeployment(); err != nil {
		return err
	}
	if err := c.createOnosConfigProxyConfigMap(); err != nil {
		return err
	}
	if err := c.createOnosConfigProxyDeployment(); err != nil {
		return err
	}
	if err := c.createOnosConfigProxyService(); err != nil {
		return err
	}
	if err := c.awaitOnosConfigDeploymentReady(); err != nil {
		return err
	}
	if err := c.awaitOnosConfigProxyDeploymentReady(); err != nil {
		return err
	}
	return nil
}

// createOnosConfigConfigMap creates a ConfigMap for the onos-config Deployment
func (c *ClusterController) createOnosConfigConfigMap() error {
	config, err := c.config.load()
	if err != nil {
		return err
	}

	// Serialize the change store configuration
	changeStore, err := json.Marshal(config["changeStore"])
	if err != nil {
		return err
	}

	// Serialize the network store configuration
	networkStore, err := json.Marshal(config["networkStore"])
	if err != nil {
		return err
	}

	// Serialize the device store configuration
	deviceStore, err := json.Marshal(config["deviceStore"])
	if err != nil {
		return err
	}

	// Serialize the config store configuration
	configStore, err := json.Marshal(config["configStore"])
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-config",
			Namespace: c.clusterID,
		},
		Data: map[string]string{
			"changeStore.json":  string(changeStore),
			"configStore.json":  string(configStore),
			"deviceStore.json":  string(deviceStore),
			"networkStore.json": string(networkStore),
		},
	}
	_, err = c.kubeclient.CoreV1().ConfigMaps(c.clusterID).Create(cm)
	return err
}

// createModelPluginString creates model plugin path based on a device type, version, and image tag
func (c *ClusterController) createModelPluginString(deviceType string, version string, debug bool) string {
	var sb strings.Builder
	sb.WriteString("-modelPlugin=/usr/local/lib/")
	sb.WriteString(deviceType)
	if debug {
		sb.WriteString("-debug.so.")
		sb.WriteString(version)
	} else {
		sb.WriteString(".so.")
		sb.WriteString(version)
	}

	return sb.String()
}

// createOnosConfigDeployment creates an onos-config Deployment
func (c *ClusterController) createOnosConfigDeployment() error {
	nodes := int32(c.config.ConfigNodes)
	zero := int64(0)

	testDevModelPluginV1 := ""
	testDevModelPluginV2 := ""
	testDevSimModelPluginV1 := ""
	testStratumModelPluginV2 := ""

	if c.config.ImageTags["config"] == string(Debug) {
		testDevModelPluginV1 = c.createModelPluginString("testdevice", "1.0.0", true)
		testDevModelPluginV2 = c.createModelPluginString("testdevice", "2.0.0", true)
		testDevSimModelPluginV1 = c.createModelPluginString("devicesim", "1.0.0", true)
		testStratumModelPluginV2 = c.createModelPluginString("stratum", "1.0.0", true)

	} else {
		testDevModelPluginV1 = c.createModelPluginString("testdevice", "1.0.0", false)
		testDevModelPluginV2 = c.createModelPluginString("testdevice", "2.0.0", false)
		testDevSimModelPluginV1 = c.createModelPluginString("devicesim", "1.0.0", false)
		testStratumModelPluginV2 = c.createModelPluginString("stratum", "1.0.0", false)

	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-config",
			Namespace: c.clusterID,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nodes,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":      "onos",
					"type":     "config",
					"resource": "onos-config",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":      "onos",
						"type":     "config",
						"resource": "onos-config",
					},
					Annotations: map[string]string{
						"seccomp.security.alpha.kubernetes.io/pod": "unconfined",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "onos-config",
							Image:           c.imageName("onosproject/onos-config", c.config.ImageTags["config"]),
							ImagePullPolicy: c.config.PullPolicy,
							Env: []corev1.EnvVar{
								{
									Name:  "ATOMIX_CONTROLLER",
									Value: fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", c.clusterID),
								},
								{
									Name:  "ATOMIX_APP",
									Value: "onos-config",
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
								"-caPath=/etc/onos-config/certs/onf.cacrt",
								"-keyPath=/etc/onos-config/certs/onos-config.key",
								"-certPath=/etc/onos-config/certs/onos-config.crt",
								"-configStore=/etc/onos-config/configs/configStore.json",
								"-changeStore=/etc/onos-config/configs/changeStore.json",
								"-networkStore=/etc/onos-config/configs/networkStore.json",
								testDevModelPluginV1,
								testDevModelPluginV2,
								testDevSimModelPluginV1,
								testStratumModelPluginV2,
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
									Name:      "config",
									MountPath: "/etc/onos-config/configs",
									ReadOnly:  true,
								},
								{
									Name:      "secret",
									MountPath: "/etc/onos-config/certs",
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
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "onos-config",
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

// createOnosConfigService creates a Service to expose the onos-config Deployment to other pods
func (c *ClusterController) createOnosConfigService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-config",
			Namespace: c.clusterID,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":      "onos",
				"type":     "config",
				"resource": "onos-config",
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

// awaitOnosConfigDeploymentReady waits for the onos-config pods to complete startup
func (c *ClusterController) awaitOnosConfigDeploymentReady() error {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": "onos", "resource": "onos-config"}}
	unblocked := make(map[string]bool)
	for {
		// Get a list of the pods that match the deployment
		pods, err := c.kubeclient.CoreV1().Pods(c.clusterID).List(metav1.ListOptions{
			LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
		})
		if err != nil {
			return err
		}

		// Iterate through the pods in the deployment and unblock the debugger
		for _, pod := range pods.Items {
			if _, ok := unblocked[pod.Name]; !ok && len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].State.Running != nil {
				unblocked[pod.Name] = true
			}
		}

		// Get the onos-config deployment
		dep, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Get("onos-config", metav1.GetOptions{})
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

// createOnosConfigProxyConfigMap creates a ConfigMap for the onos-config-envoy Deployment
func (c *ClusterController) createOnosConfigProxyConfigMap() error {
	configPath := filepath.Join(filepath.Join(configsPath, "envoy"), "envoy-config.yaml")
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-config-envoy",
			Namespace: c.clusterID,
		},
		BinaryData: map[string][]byte{
			"envoy-config.yaml": data,
		},
	}
	_, err = c.kubeclient.CoreV1().ConfigMaps(c.clusterID).Create(cm)
	return err
}

// createOnosConfigProxyDeployment creates an onos-config Envoy proxy
func (c *ClusterController) createOnosConfigProxyDeployment() error {
	nodes := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-config-envoy",
			Namespace: c.clusterID,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nodes,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  "onos",
					"type": "config-envoy",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":      "onos",
						"type":     "config-envoy",
						"resource": "onos-config-envoy",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "onos-config-envoy",
							Image:           "envoyproxy/envoy-alpine:latest",
							ImagePullPolicy: c.config.PullPolicy,
							Command: []string{
								"/usr/local/bin/envoy",
								"-c",
								"/etc/envoy-proxy/config/envoy-config.yaml",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "envoy",
									ContainerPort: 8080,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/envoy-proxy/config",
									ReadOnly:  true,
								},
								{
									Name:      "secret",
									MountPath: "/etc/envoy-proxy/certs",
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
										Name: "onos-config-envoy",
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
	_, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Create(deployment)
	return err
}

// createOnosConfigProxyService creates an onos-config Envoy proxy service
func (c *ClusterController) createOnosConfigProxyService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-config-envoy",
			Namespace: c.clusterID,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":      "onos",
				"type":     "config-envoy",
				"resource": "onos-config-envoy",
			},
			Ports: []corev1.ServicePort{
				{
					Name: "envoy",
					Port: 8080,
				},
			},
		},
	}
	_, err := c.kubeclient.CoreV1().Services(c.clusterID).Create(service)
	return err
}

// awaitOnosConfigProxyDeploymentReady waits for the onos-config proxy pods to complete startup
func (c *ClusterController) awaitOnosConfigProxyDeploymentReady() error {
	for {
		// Get the onos-config deployment
		dep, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Get("onos-config-envoy", metav1.GetOptions{})
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

// GetOnosConfigNodes returns a list of all onos-config nodes running in the cluster
func (c *ClusterController) GetOnosConfigNodes() ([]NodeInfo, error) {
	configLabelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": "onos", "resource": "onos-config"}}

	pods, err := c.kubeclient.CoreV1().Pods(c.clusterID).List(metav1.ListOptions{
		LabelSelector: labels.Set(configLabelSelector.MatchLabels).String(),
	})
	if err != nil {
		return nil, err
	}

	onosConfigNodes := make([]NodeInfo, len(pods.Items))
	for i, pod := range pods.Items {
		var status NodeStatus
		if pod.Status.Phase == corev1.PodRunning {
			status = NodeRunning
		} else if pod.Status.Phase == corev1.PodFailed {
			status = NodeFailed
		}
		onosConfigNodes[i] = NodeInfo{
			ID:     pod.Name,
			Status: status,
			Type:   OnosConfig,
		}
	}

	return onosConfigNodes, nil
}
