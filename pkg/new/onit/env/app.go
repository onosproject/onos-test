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

package env

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/new/onit/setup"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"
)

func newAppSetup(name string, testEnv *testEnv) setup.AppSetup {
	setup := &appSetup{
		serviceSetup: &serviceSetup{
			testEnv: testEnv,
		},
		name: name,
	}
	setup.serviceSetup.setup = setup
	return setup
}

// App provides the environment for an app
type App interface {
	Service
}

var _ App = &app{}

// app is an implementation of the App interface
type app struct {
	*service
}

func (e *app) Remove() {
	if err := e.teardownApp(); err != nil {
		panic(err)
	}
}

// teardownApp tears down a app by name
func (e *app) teardownApp() error {
	var err error
	if e := e.deleteAppDeployment(); e != nil {
		err = e
	}
	if e := e.deleteAppService(); e != nil {
		err = e
	}
	if e := e.deleteAppConfigMap(); e != nil {
		err = e
	}
	return err
}

// deleteAppConfigMap deletes an app ConfigMap by name
func (e *app) deleteAppConfigMap() error {
	return e.kubeClient.CoreV1().ConfigMaps(e.namespace).Delete(e.name, &metav1.DeleteOptions{})
}

// deleteAppPod deletes an app Pod by name
func (e *app) deleteAppDeployment() error {
	return e.kubeClient.AppsV1().Deployments(e.namespace).Delete(e.name, &metav1.DeleteOptions{})
}

// deleteAppService deletes an app Service by name
func (e *app) deleteAppService() error {
	return e.kubeClient.CoreV1().Services(e.namespace).Delete(e.name, &metav1.DeleteOptions{})
}

var _ setup.AppSetup = &appSetup{}

// appSetup is an implementation of the AppSetup interface
type appSetup struct {
	*serviceSetup
	name  string
	nodes int
}

func (s *appSetup) Nodes(nodes int) setup.AppSetup {
	s.nodes = nodes
	return s
}

func (s *appSetup) Using() setup.ServiceSetup {
	return s
}

func (s *appSetup) Setup() error {
	if err := s.createService(); err != nil {
		return err
	}
	if err := s.createDeployment(); err != nil {
		return err
	}
	if err := s.awaitDeploymentReady(); err != nil {
		return err
	}
	return nil
}

// createDeployment creates an app Deployment
func (s *appSetup) createDeployment() error {
	nodes := int32(1)
	zero := int64(0)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
			Labels: map[string]string{
				"type": "app",
				"app":  s.name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nodes,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"type": "app",
					"app":  s.name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"type": "app",
						"app":  s.name,
					},
					Annotations: map[string]string{
						"seccomp.security.alpha.kubernetes.io/pod": "unconfined",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            s.name,
							Image:           s.image,
							ImagePullPolicy: s.pullPolicy,
							Env: []corev1.EnvVar{
								{
									Name:  "ATOMIX_CONTROLLER",
									Value: fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", s.namespace),
								},
								{
									Name:  "ATOMIX_APP",
									Value: "test",
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

// createService creates an app service
func (s *appSetup) createService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"type": "app",
				"app":  s.name,
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

// awaitDeploymentReady waits for the app pods to complete startup
func (s *appSetup) awaitDeploymentReady() error {
	for {
		// Get the app deployment
		dep, err := s.kubeClient.AppsV1().Deployments(s.namespace).Get(s.name, metav1.GetOptions{})
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
