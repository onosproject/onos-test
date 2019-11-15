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
	DeploymentEnv

	// Address returns the service address
	Address() string

	// Execute executes the given command and returns the output
	Execute(command ...string) ([]string, int, error)

	// Credentials returns the service credentials
	Credentials() *tls.Config

	// Connect connects to the service
	Connect() (*grpc.ClientConn, error)
}

// clusterServiceEnv is an implementation of the Service interface
type clusterServiceEnv struct {
	*clusterDeploymentEnv
	service *cluster.Service
}

func (e *clusterServiceEnv) Address() string {
	return e.service.Address()
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
