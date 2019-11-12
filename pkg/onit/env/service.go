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
	"crypto/tls"
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"google.golang.org/grpc"
)

// Service is a base interface for service environments
type Service interface {
	// Address returns the service address
	Address() string

	// Name is the name of the service
	Name() string

	// Nodes returns the service nodes
	Nodes() []Node

	// Node returns a specific node environment
	Node(name string) Node

	// AwaitReady waits for all nodes in the service to become ready
	AwaitReady() error

	// Execute executes the given command and returns the output
	Execute(command ...string) ([]string, int, error)

	// Credentials returns the service credentials
	Credentials() *tls.Config

	// Connect connects to the service
	Connect() (*grpc.ClientConn, error)
}

// ServiceSetup is a base interface for services that can be set up
type ServiceSetup interface {
	Service
	setup.Setup
}

// clusterService is an implementation of the Service interface
type clusterService struct {
	service *cluster.Service
}

func (e *clusterService) Name() string {
	return e.service.Name()
}

func (e *clusterService) Address() string {
	return e.service.Address()
}

func (e *clusterService) Nodes() []Node {
	clusterNodes := e.service.Nodes().List()
	nodes := make([]Node, len(clusterNodes))
	for i, node := range clusterNodes {
		nodes[i] = e.Node(node.Name())
	}
	return nodes
}

func (e *clusterService) Node(name string) Node {
	return &clusterNode{
		e.service.Nodes().Get(name),
	}
}

func (e *clusterService) AwaitReady() error {
	return e.service.AwaitReady()
}

func (e *clusterService) Execute(command ...string) ([]string, int, error) {
	return e.service.Execute(command...)
}

func (e *clusterService) Credentials() *tls.Config {
	config, err := e.service.Credentials()
	if err != nil {
		panic(err)
	}
	return config
}

func (e *clusterService) Connect() (*grpc.ClientConn, error) {
	return e.service.Connect()
}
