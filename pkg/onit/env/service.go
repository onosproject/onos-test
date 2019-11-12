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
	"google.golang.org/grpc"
)

// ServiceEnv is a base interface for service environments
type ServiceEnv interface {
	// Address returns the service address
	Address() string

	// Name is the name of the service
	Name() string

	// Nodes returns the service nodes
	Nodes() []NodeEnv

	// Node returns a specific node environment
	Node(name string) NodeEnv

	// AwaitReady waits for all nodes in the service to become ready
	AwaitReady() error

	// Execute executes the given command and returns the output
	Execute(command ...string) ([]string, int, error)

	// Credentials returns the service credentials
	Credentials() *tls.Config

	// Connect connects to the service
	Connect() (*grpc.ClientConn, error)
}

// clusterServiceEnv is an implementation of the Service interface
type clusterServiceEnv struct {
	service *cluster.Service
}

func (e *clusterServiceEnv) Name() string {
	return e.service.Name()
}

func (e *clusterServiceEnv) Address() string {
	return e.service.Address()
}

func (e *clusterServiceEnv) Nodes() []NodeEnv {
	clusterNodes := e.service.Nodes().List()
	nodes := make([]NodeEnv, len(clusterNodes))
	for i, node := range clusterNodes {
		nodes[i] = e.Node(node.Name())
	}
	return nodes
}

func (e *clusterServiceEnv) Node(name string) NodeEnv {
	return &clusterNodeEnv{
		e.service.Nodes().Get(name),
	}
}

func (e *clusterServiceEnv) AwaitReady() error {
	return e.service.AwaitReady()
}

func (e *clusterServiceEnv) Execute(command ...string) ([]string, int, error) {
	return e.service.Execute(command...)
}

func (e *clusterServiceEnv) Credentials() *tls.Config {
	config, err := e.service.Credentials()
	if err != nil {
		panic(err)
	}
	return config
}

func (e *clusterServiceEnv) Connect() (*grpc.ClientConn, error) {
	return e.service.Connect()
}
