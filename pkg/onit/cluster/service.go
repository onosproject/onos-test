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
	"fmt"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"path"
	"strings"
	"time"
)

func newService(name string, ports []Port, labels map[string]string, image string, secrets map[string]string, args []string, client *client) *Service {
	return &Service{
		Deployment: newDeployment(name, labels, image, client),
		replicas:   GetArg(name, "replicas").Int(1),
		ports:      ports,
		secrets:    secrets,
		env:        make(map[string]string),
		args:       args,
	}
}

// Port is a service port
type Port struct {
	Name string
	Port int
}

// Service is the base type for multi-node services
type Service struct {
	*Deployment
	ports      []Port
	replicas   int
	debug      bool
	user       *int
	privileged bool
	secrets    map[string]string
	env        map[string]string
	args       []string
}

// SetName sets the service name
func (s *Service) SetName(name string) {
	s.name = name
}

// Ports returns the service ports
func (s *Service) Ports() []Port {
	return s.ports
}

// SetPorts sets the service ports
func (s *Service) SetPorts(ports []Port) {
	s.ports = ports
}

// AddPort adds a port to the service
func (s *Service) AddPort(name string, port int) {
	s.ports = append(s.ports, Port{
		Name: name,
		Port: port,
	})
}

// Debug returns whether debug is enabled
func (s *Service) Debug() bool {
	return s.debug
}

// SetDebug sets whether debug is enabled
func (s *Service) SetDebug(debug bool) {
	s.debug = debug
}

// getPort gets a port by name
func (s *Service) getPortByName(name string) *Port {
	for _, port := range s.ports {
		if port.Name == name {
			return &port
		}
	}
	return nil
}

// Address returns the service address for the given port
func (s *Service) Address(port string) string {
	info := s.getPortByName(port)
	if info == nil {
		panic(fmt.Errorf("unknown port %s", port))
	}
	return fmt.Sprintf("%s:%d", s.name, info.Port)
}

// Replicas returns the number of nodes in the service
func (s *Service) Replicas() int {
	return s.replicas
}

// SetReplicas sets the number of nodes in the service
func (s *Service) SetReplicas(replicas int) {
	s.replicas = replicas
}

// User returns the user with which to run the service
func (s *Service) User() *int {
	return s.user
}

// SetUser sets the user with which to run the service
func (s *Service) SetUser(user int) {
	s.user = &user
}

// Privileged returns whether to run the service in privileged mode
func (s *Service) Privileged() bool {
	return s.privileged
}

// SetPrivileged sets whether to run the service in privileged mode
func (s *Service) SetPrivileged(privileged bool) {
	s.privileged = privileged
}

// Secrets returns the service secrets
func (s *Service) Secrets() map[string]string {
	return s.secrets
}

// SetSecrets sets the service secrets
func (s *Service) SetSecrets(secrets map[string]string) {
	s.secrets = secrets
}

// AddSecret adds a secret to the service
func (s *Service) AddSecret(name, secret string) {
	s.secrets[name] = secret
}

// Env returns the service environment
func (s *Service) Env() map[string]string {
	return s.env
}

// SetEnv sets the service environment
func (s *Service) SetEnv(env map[string]string) {
	s.env = env
}

// AddEnv adds an environment variable
func (s *Service) AddEnv(name, value string) {
	s.env[name] = value
}

// Args returns the service arguments
func (s *Service) Args() []string {
	return s.args
}

// SetArgs sets the service arguments
func (s *Service) SetArgs(args ...string) {
	s.args = args
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

func getKey(key string) string {
	return strings.ReplaceAll(path.Base(key), "/", "-")
}

// createSecret creates the service Secret
func (s *Service) createSecret() error {
	if len(s.secrets) == 0 {
		return nil
	}

	secrets := make(map[string]string)
	for key, value := range s.secrets {
		secrets[getKey(key)] = value
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name(),
			Namespace: s.namespace,
			Labels:    s.labels,
		},
		StringData: secrets,
	}
	_, err := s.kubeClient.CoreV1().Secrets(s.namespace).Create(secret)
	return err
}

// createService creates a Service to expose the Deployment to other pods
func (s *Service) createService() error {
	ports := make([]corev1.ServicePort, len(s.ports))
	for i, port := range s.ports {
		ports[i] = corev1.ServicePort{
			Name: port.Name,
			Port: int32(port.Port),
		}
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name(),
			Namespace: s.namespace,
			Labels:    s.labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: s.labels,
			Ports:    ports,
		},
	}
	_, err := s.kubeClient.CoreV1().Services(s.namespace).Create(service)
	return err
}

// createDeployment creates a Deployment
func (s *Service) createDeployment() error {
	nodes := int32(s.Replicas())

	// Expose all provided ports
	ports := make([]corev1.ContainerPort, len(s.ports))
	for i, port := range s.ports {
		ports[i] = corev1.ContainerPort{
			Name:          port.Name,
			ContainerPort: int32(port.Port),
		}
	}

	// If the debug port is set, assume debugging is enabled and set the security contexts appropriately
	var securityContext *corev1.SecurityContext
	var podSecurityContext *corev1.PodSecurityContext

	if s.Privileged() {
		privileged := true
		securityContext = &corev1.SecurityContext{
			Privileged: &privileged,
		}
	}

	if s.user != nil {
		user := int64(*s.user)
		podSecurityContext = &corev1.PodSecurityContext{
			RunAsUser: &user,
		}
	}

	if s.Debug() {
		zero := int64(0)
		podSecurityContext = &corev1.PodSecurityContext{
			RunAsUser: &zero,
		}
		securityContext = &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{
					"SYS_PTRACE",
				},
			},
		}
	}

	if s.Privileged() {
		privileged := true
		securityContext = &corev1.SecurityContext{
			Privileged: &privileged,
		}
	}

	if s.user != nil {
		user := int64(*s.user)
		podSecurityContext = &corev1.PodSecurityContext{
			RunAsUser: &user,
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

		volumeMounts = make([]corev1.VolumeMount, 0, len(s.secrets))
		for filepath := range s.secrets {
			filename := path.Dir(filepath)
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      filename,
				MountPath: filepath,
				SubPath:   getKey(filepath),
				ReadOnly:  true,
			})
		}
		volumeMounts = []corev1.VolumeMount{
			{
				Name:      "secret",
				MountPath: "/certs",
				ReadOnly:  true,
			},
		}
	}

	var readinessProbe *corev1.Probe
	var livenessProbe *corev1.Probe
	if len(s.ports) > 0 {
		port := s.ports[0]
		readinessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.FromInt(port.Port),
				},
			},
			InitialDelaySeconds: 5,
			PeriodSeconds:       10,
		}
		livenessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.FromInt(port.Port),
				},
			},
			InitialDelaySeconds: 15,
			PeriodSeconds:       20,
		}
	}

	env := []corev1.EnvVar{
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
			Value: "database",
		},
	}
	for name, value := range s.env {
		env = append(env, corev1.EnvVar{
			Name:  name,
			Value: value,
		})
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
							Env:             env,
							Args:            s.args,
							Ports:           ports,
							ReadinessProbe:  readinessProbe,
							LivenessProbe:   livenessProbe,
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

// Credentials returns the TLS credentials
func (s *Service) Credentials() (*tls.Config, error) {
	return getClientCredentials()
}

// Connect creates a gRPC client connection to the service
func (s *Service) Connect(port string) (*grpc.ClientConn, error) {
	tlsConfig, err := s.Credentials()
	if err != nil {
		return nil, err
	}
	return grpc.Dial(s.Address(port), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}
