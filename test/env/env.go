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
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	atomix "github.com/atomix/atomix-go-client/pkg/client"
	"github.com/onosproject/onos-config/pkg/northbound/admin"
	"github.com/openconfig/gnmi/client"
	gnmi "github.com/openconfig/gnmi/client/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	// TestDevicesEnv : environment variable name for devices
	TestDevicesEnv = "ONOS_CONFIG_TEST_DEVICES"
)

const (
	clientKeyPath = "/etc/onos-config/certs/client1.key"
	clientCrtPath = "/etc/onos-config/certs/client1.crt"
	caCertPath    = "/etc/onos-config/certs/onf.cacrt"
	configAddress = "onos-config:5150"
	topoAddress   = "onos-topo:5150"
)

// ExecuteCLI executes an onos CLI command and returns the output and exit code
func ExecuteCLI(command ...string) ([]string, int) {
	nodes := GetCLINodes()
	return ExecuteCommand(nodes[0], append([]string{"/bin/bash", "-c"}, command...)...)
}

// GetCredentials returns gNMI client credentials for the test environment
func GetCredentials() (*tls.Config, error) {
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	if !certPool.AppendCertsFromPEM(ca) {
		return nil, errors.New("failed to append CA certificates")
	}

	cert, err := tls.LoadX509KeyPair(clientCrtPath, clientKeyPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs:            certPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}, nil
}

// GetDestination returns a gNMI client destination for the test environment
func GetDestination(target string) (client.Destination, error) {
	tlsConfig, err := GetCredentials()
	if err != nil {
		return client.Destination{}, err
	}
	return client.Destination{
		Addrs:   []string{configAddress},
		Target:  target,
		TLS:     tlsConfig,
		Timeout: 10 * time.Second,
	}, nil
}

// GetDestinationForDevice returns a gNMI client destination for the test environment
func GetDestinationForDevice(addr string, target string) (client.Destination, error) {
	tlsConfig, err := GetCredentials()
	if err != nil {
		return client.Destination{}, err
	}
	return client.Destination{
		Addrs:   []string{addr},
		Target:  target,
		TLS:     tlsConfig,
		Timeout: 10 * time.Second,
	}, nil
}

// NewGnmiClient returns a new gNMI client for the test environment
func NewGnmiClient(ctx context.Context, target string) (client.Impl, error) {
	dest, err := GetDestination(target)
	if err != nil {
		return nil, err
	}
	return gnmi.New(ctx, dest)
}

// NewGnmiClientForDevice returns a new gNMI client for the test environment
func NewGnmiClientForDevice(ctx context.Context, address string, target string) (client.Impl, error) {
	dest, err := GetDestinationForDevice(address, target)
	if err != nil {
		return nil, err
	}
	insecureConnection, insecureConnectionError := getInsecureConn(address)
	if insecureConnectionError != nil {
		return nil, insecureConnectionError
	}
	return gnmi.NewFromConn(ctx, insecureConnection, dest)
}

// GetDevices returns a slice of device names for the test environment
func GetDevices() []string {
	devices := os.Getenv(TestDevicesEnv)
	return strings.Split(devices, ",")
}

func handleCertArgs() ([]grpc.DialOption, error) {
	var opts = []grpc.DialOption{}
	var cert tls.Certificate
	var err error

	// Load default Certificates
	cert, err = tls.LoadX509KeyPair(clientCrtPath, clientKeyPath)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))

	return opts, nil
}

func handleInsecureArgs() ([]grpc.DialOption, error) {
	var opts = []grpc.DialOption{}

	opts = append(opts, grpc.WithInsecure())

	return opts, nil
}

// getConn gets a gRPC connection to the given address
func getConn(address string) (*grpc.ClientConn, error) {
	opts, err := handleCertArgs()
	if err != nil {
		return nil, err
	}
	return grpc.Dial(address, opts...)
}

// getInsecureConn gets a gRPC connection to the given address
func getInsecureConn(address string) (*grpc.ClientConn, error) {
	opts, err := handleInsecureArgs()
	if err != nil {
		return nil, err
	}
	return grpc.Dial(address, opts...)
}

// GetTopoConn gets a gRPC connection to the topology service
func GetTopoConn() (*grpc.ClientConn, error) {
	return getConn(topoAddress)
}

// GetConfigConn gets a gRPC connection to the config service
func GetConfigConn() (*grpc.ClientConn, error) {
	return getConn(configAddress)
}

// GetAdminClient returns a client that can be used for the admin APIs
func GetAdminClient() (*grpc.ClientConn, admin.ConfigAdminServiceClient) {
	opts, err := handleCertArgs()
	if err != nil {
		fmt.Printf("Error loading cert %s", err)
	}
	conn, err := grpc.Dial(configAddress, opts...)
	if err != nil {
		panic(err)
	}
	return conn, admin.NewConfigAdminServiceClient(conn)
}

// NewAtomixClient returns an Atomix client from the environment
func NewAtomixClient(test string) (*atomix.Client, error) {
	opts := []atomix.Option{
		atomix.WithNamespace(os.Getenv("ATOMIX_NAMESPACE")),
		atomix.WithApplication(test),
	}
	return atomix.NewClient(os.Getenv("ATOMIX_CONTROLLER"), opts...)
}
