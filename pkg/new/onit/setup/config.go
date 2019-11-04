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

package setup

import (
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config is an interface for setting up config nodes
type Config interface {
	ServiceType
	concurrentSetup

	// Nodes sets the number of nodes to deploy
	Nodes(nodes int) Config
}

var _ Config = &config{}

// config is an implementation of the Config interface
type config struct {
	*serviceType
	nodes int
}

func (s *config) Nodes(nodes int) Config {
	s.nodes = nodes
	return s
}

func (s *config) create() error {
	if err := s.createOnosConfigService(); err != nil {
		return err
	}
	if err := s.createOnosConfigDeployment(); err != nil {
		return err
	}
	return nil
}

// createModelPluginString creates model plugin path based on a device type, version, and image tag
func (s *config) createModelPluginString(deviceType string, version string, debug bool) string {
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
func (s *config) createOnosConfigDeployment() error {
	nodes := int32(s.nodes)
	zero := int64(0)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-config",
			Namespace: s.namespace,
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
							Image:           s.image,
							ImagePullPolicy: s.pullPolicy,
							Env: []corev1.EnvVar{
								{
									Name:  "ATOMIX_CONTROLLER",
									Value: fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", s.namespace),
								},
								{
									Name:  "ATOMIX_APP",
									Value: "onos-config",
								},
								{
									Name:  "ATOMIX_NAMESPACE",
									Value: s.namespace,
								},
								{
									Name:  "ATOMIX_RAFT_GROUP",
									Value: "raft",
								},
								{
									Name: "NODE_ID",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.name",
										},
									},
								},
							},
							Args: []string{
								"-caPath=/etc/onos-config/certs/onf.cacrt",
								"-keyPath=/etc/onos-config/certs/onos-config.key",
								"-certPath=/etc/onos-config/certs/onos-config.crt",
								"-modelPlugin=/usr/local/lib/testdevice.so.1.0.0",
								"-modelPlugin=/usr/local/lib/testdevice.so.2.0.0",
								"-modelPlugin=/usr/local/lib/devicesim.so.1.0.0",
								"-modelPlugin=/usr/local/lib/stratum.so.1.0.0",
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
									SecretName: s.namespace,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := s.kubeClient.AppsV1().Deployments(s.namespace).Create(dep)
	return err
}

// createOnosConfigService creates a Service to expose the onos-config Deployment to other pods
func (s *config) createOnosConfigService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-config",
			Namespace: s.namespace,
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
	_, err := s.kubeClient.CoreV1().Services(s.namespace).Create(service)
	return err
}

// waitForStart waits for the onos-config pods to complete startup
func (s *config) waitForStart() error {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": "onos", "resource": "onos-config"}}
	unblocked := make(map[string]bool)
	for {
		// Get a list of the pods that match the deployment
		pods, err := s.kubeClient.CoreV1().Pods(s.namespace).List(metav1.ListOptions{
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
		dep, err := s.kubeClient.AppsV1().Deployments(s.namespace).Get("onos-config", metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Return once the all replicas in the deployment are ready
		if int(dep.Status.ReadyReplicas) == s.nodes {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}
