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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	corev1 "k8s.io/api/core/v1"
)

func newService(name string, port int, labels map[string]string, image string, client *client) *Service {
	return &Service{
		client: client,
		name:   name,
		port:   port,
		labels: labels,
		image:  image,
	}
}

// Service is the base type for multi-node services
type Service struct {
	*client
	name       string
	port       int
	replicas   int
	labels     map[string]string
	image      string
	pullPolicy corev1.PullPolicy
}

// Name returns the name of the service
func (s *Service) Name() string {
	return s.name
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

// Address returns the service address
func (s *Service) Address() string {
	return fmt.Sprintf("%s:%d", s.name, s.port)
}

// Nodes returns the collection of nodes in the service
func (s *Service) Nodes() *Nodes {
	return newNodes(s.port, s.labels, s.client)
}

// Replicas returns the number of nodes in the service
func (s *Service) Replicas() int {
	return s.replicas
}

// SetReplicas sets the number of nodes in the service
func (s *Service) SetReplicas(replicas int) {
	s.replicas = replicas
}

// Image returns the image for the service
func (s *Service) Image() string {
	return s.image
}

// SetImage sets the image for the service
func (s *Service) SetImage(image string) {
	s.image = image
}

// PullPolicy returns the image pull policy for the service
func (s *Service) PullPolicy() corev1.PullPolicy {
	return s.pullPolicy
}

// SetPullPolicy sets the image pull policy for the service
func (s *Service) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	s.pullPolicy = pullPolicy
}

// AwaitReady waits for the service to become ready
func (s *Service) AwaitReady() error {
	return s.Nodes().AwaitReady()
}

// Execute executes the given command on one of the service nodes
func (s *Service) Execute(command ...string) ([]string, int, error) {
	nodes := s.Nodes().List()
	if len(nodes) == 0 {
		return nil, 0, errors.New("no service nodes found")
	}
	return nodes[0].Execute(command...)
}

// Credentials returns the TLS credentials
func (s *Service) Credentials() (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}, nil
}

// Connect creates a gRPC client connection to the service
func (s *Service) Connect() (*grpc.ClientConn, error) {
	tlsConfig, err := s.Credentials()
	if err != nil {
		return nil, err
	}
	return grpc.Dial(s.Address(), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}
