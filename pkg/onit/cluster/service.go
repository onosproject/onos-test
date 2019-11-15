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
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"
)

func newService(name string, port int, labels map[string]string, image string, secrets map[string]string, args []string, client *client) *Service {
	return &Service{
		Deployment: newDeployment(name, labels, image, client),
		port:       port,
		secrets:    secrets,
		args:       args,
	}
}

// Service is the base type for multi-node services
type Service struct {
	*Deployment
	port      int
	debugPort int
	replicas  int
	secrets   map[string]string
	args      []string
}

// SetName sets the service name
func (s *Service) SetName(name string) {
	s.name = name
}

// Port returns the service port
func (s *Service) Port() int {
	return s.port
}

// SetPort sets the service port
func (s *Service) SetPort(port int) {
	s.port = port
}

// DebugPort returns the service debug port
func (s *Service) DebugPort() int {
	return s.debugPort
}

// SetDebugPort sets the service debug port
func (s *Service) SetDebugPort(port int) {
	s.debugPort = port
}

// Address returns the service address
func (s *Service) Address() string {
	return fmt.Sprintf("%s:%d", s.name, s.port)
}

// Replicas returns the number of nodes in the service
func (s *Service) Replicas() int {
	return s.replicas
}

// SetReplicas sets the number of nodes in the service
func (s *Service) SetReplicas(replicas int) {
	s.replicas = replicas
}

// Setup sets up the service
func (s *Service) Setup() error {
	if s.replicas == 0 {
		return nil
	}

	step := logging.NewStep(s.namespace, "Setup %s", s.Name())
	step.Start()
	step.Logf("Creating %s Service", s.Name())
	if err := s.createService(); err != nil {
		step.Fail(err)
		return err
	}
	step.Logf("Creating %s Secret", s.Name())
	if err := s.createSecret(); err != nil {
		step.Fail(err)
		return err
	}
	step.Logf("Creating %s Deployment", s.Name())
	if err := s.createDeployment(); err != nil {
		step.Fail(err)
		return err
	}
	step.Logf("Waiting for %s to become ready", s.Name())
	if err := s.AwaitReady(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// createSecret creates the service Secret
func (s *Service) createSecret() error {
	if len(s.secrets) == 0 {
		return nil
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name(),
			Namespace: s.namespace,
			Labels:    s.labels,
		},
		StringData: s.secrets,
	}
	_, err := s.kubeClient.CoreV1().Secrets(s.namespace).Create(secret)
	return err
}

// createService creates a Service to expose the Deployment to other pods
func (s *Service) createService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name(),
			Namespace: s.namespace,
			Labels:    s.labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: s.labels,
			Ports: []corev1.ServicePort{
				{
					Name: "grpc",
					Port: int32(s.Port()),
				},
			},
		},
	}
	_, err := s.kubeClient.CoreV1().Services(s.namespace).Create(service)
	return err
}

// createDeployment creates a Deployment
func (s *Service) createDeployment() error {
	nodes := int32(s.Replicas())
	zero := int64(0)

	// Default to exposing only a single port
	ports := []corev1.ContainerPort{
		{
			Name:          "grpc",
			ContainerPort: int32(s.Port()),
		},
	}

	// If the debug port is set, assume debugging is enabled and set the security contexts appropriately
	var securityContext *corev1.SecurityContext
	var podSecurityContext *corev1.PodSecurityContext
	if s.DebugPort() != 0 {
		ports = append(ports, corev1.ContainerPort{
			Name:          "debug",
			ContainerPort: int32(s.DebugPort()),
			Protocol:      corev1.ProtocolTCP,
		})
		securityContext = &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{
					"SYS_PTRACE",
				},
			},
		}
		podSecurityContext = &corev1.PodSecurityContext{
			RunAsUser: &zero,
		}
	}

	// Mount secrets if necessary
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	if len(s.secrets) > 0 {
		volumes = []corev1.Volume{
			{
				Name: "secret",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: s.Name(),
					},
				},
			},
		}
		volumeMounts = []corev1.VolumeMount{
			{
				Name:      "secret",
				MountPath: "/certs",
				ReadOnly:  true,
			},
		}
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name(),
			Namespace: s.namespace,
			Labels:    s.labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nodes,
			Selector: &metav1.LabelSelector{
				MatchLabels: s.labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: s.labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            s.Name(),
							Image:           s.Image(),
							ImagePullPolicy: s.PullPolicy(),
							Env: []corev1.EnvVar{
								{
									Name:  "ATOMIX_CONTROLLER",
									Value: fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", s.namespace),
								},
								{
									Name:  "ATOMIX_APP",
									Value: s.Name(),
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
							Args:  s.args,
							Ports: ports,
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(s.Port()),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
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
							VolumeMounts:    volumeMounts,
							SecurityContext: securityContext,
						},
					},
					SecurityContext: podSecurityContext,
					Volumes:         volumes,
				},
			},
		},
	}
	_, err := s.kubeClient.AppsV1().Deployments(s.namespace).Create(dep)
	return err
}

// AwaitReady waits for the service to complete startup
func (s *Service) AwaitReady() error {
	if s.replicas == 0 {
		return nil
	}

	for {
		// Get the deployment
		dep, err := s.kubeClient.AppsV1().Deployments(s.namespace).Get(s.Name(), metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Return once the all replicas in the deployment are ready
		if int(dep.Status.ReadyReplicas) == s.Replicas() {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// TearDown removes the service from the cluster
func (s *Service) TearDown() error {
	var err error
	if e := s.deleteDeployment(); e != nil {
		err = e
	}
	if e := s.deleteService(); e != nil {
		err = e
	}
	if e := s.deleteSecret(); e != nil {
		err = e
	}
	return err
}

// deletePod deletes a service Deployment by name
func (s *Service) deleteDeployment() error {
	return s.kubeClient.AppsV1().Deployments(s.namespace).Delete(s.name, &metav1.DeleteOptions{})
}

// deleteService deletes a service Service by name
func (s *Service) deleteService() error {
	return s.kubeClient.CoreV1().Services(s.namespace).Delete(s.name, &metav1.DeleteOptions{})
}

// deleteSecret deletes a service Secret by name
func (s *Service) deleteSecret() error {
	_ = s.kubeClient.CoreV1().Secrets(s.namespace).Delete(s.name, &metav1.DeleteOptions{})
	return nil
}

// Execute executes the given command on one of the service nodes
func (s *Service) Execute(command ...string) ([]string, int, error) {
	nodes, err := s.Nodes()
	if err != nil {
		return nil, 0, err
	}
	if len(nodes) == 0 {
		return nil, 0, errors.New("no service nodes found")
	}
	return nodes[0].Execute(command...)
}

// Credentials returns the TLS credentials
func (s *Service) Credentials() (*tls.Config, error) {
	return getClientCredentials()
}

// Connect creates a gRPC client connection to the service
func (s *Service) Connect() (*grpc.ClientConn, error) {
	tlsConfig, err := s.Credentials()
	if err != nil {
		return nil, err
	}
	return grpc.Dial(s.Address(), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}
