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
	"github.com/onosproject/onos-test/pkg/new/util/logging"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"
)

func newTopo(client *client) *Topo {
	return &Topo{
		Service: newService("onos-topo", topoType, client),
	}
}

// Topo provides methods for managing the onos-topo service
type Topo struct {
	*Service
}

// Create creates the topology subsystem
func (s *Topo) Create() error {
	step := logging.NewStep(s.namespace, "Setup onos-topo service")
	step.Log("Create onos-topo Service")
	if err := s.createService(); err != nil {
		step.Fail(err)
		return err
	}
	step.Log("Create onos-topo Deployment")
	if err := s.createDeployment(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// createService creates a Service to expose the onos-topo Deployment to other pods
func (s *Topo) createService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-topo",
			Namespace: s.namespace,
			Labels: map[string]string{
				"type": "topo",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"type": "topo",
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

// createDeployment creates an onos-topo Deployment
func (s *Topo) createDeployment() error {
	nodes := int32(s.replicas)
	zero := int64(0)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-topo",
			Namespace: s.namespace,
			Labels: map[string]string{
				"type": "topo",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nodes,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"type": "topo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"type": "topo",
					},
					Annotations: map[string]string{
						"seccomp.security.alpha.kubernetes.io/pod": "unconfined",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "onos-topo",
							Image:           s.image,
							ImagePullPolicy: s.pullPolicy,
							Env: []corev1.EnvVar{
								{
									Name:  "ATOMIX_CONTROLLER",
									Value: fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", s.namespace),
								},
								{
									Name:  "ATOMIX_APP",
									Value: "onos-topo",
								},
								{
									Name:  "ATOMIX_NAMESPACE",
									Value: s.namespace,
								},
								{
									Name:  "ATOMIX_RAFT_GROUP",
									Value: "raft",
								},
							},
							Args: []string{
								"-caPath=/etc/onos-topo/certs/onf.cacrt",
								"-keyPath=/etc/onos-topo/certs/onos-config.key",
								"-certPath=/etc/onos-topo/certs/onos-config.crt",
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
									Name:      "secret",
									MountPath: "/etc/onos-topo/certs",
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

// AwaitReady waits for the onos-topo pods to complete startup
func (s *Topo) AwaitReady() error {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{typeLabel: topoType.name()}}
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

		// Get the onos-topo deployment
		dep, err := s.kubeClient.AppsV1().Deployments(s.namespace).Get("onos-topo", metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Return once the all replicas in the deployment are ready
		if int(dep.Status.ReadyReplicas) == s.replicas {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}
