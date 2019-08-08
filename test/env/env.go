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
	atomixclient "github.com/atomix/atomix-go-client/pkg/client"
	"github.com/onosproject/onos-config/pkg/northbound/proto"
	"github.com/openconfig/gnmi/client"
	gnmi "github.com/openconfig/gnmi/client/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	k8slabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"os"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
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

// GetNamespace returns the namespace within which the test is running
func GetNamespace() string {
	return os.Getenv("TEST_NAMESPACE")
}

// GetConfigNodes returns a list of onos-config nodes
func GetConfigNodes() []string {
	return getNodes(map[string]string{
		"app":  "onos",
		"type": "config",
	})
}

// GetTopoNodes returns a list of onos-topo nodes
func GetTopoNodes() []string {
	return getNodes(map[string]string{
		"app":  "onos",
		"type": "topo",
	})
}

// getNodes returns a list of nodes with the given labels
func getNodes(labels map[string]string) []string {
	client := mustClient()
	pods := &corev1.PodList{}
	options := &k8sclient.ListOptions{
		Namespace:     GetNamespace(),
		LabelSelector: k8slabels.SelectorFromValidatedSet(labels),
	}
	err := client.List(context.TODO(), options, pods)
	if err != nil {
		panic(err)
	}

	nodeIDs := make([]string, len(pods.Items))
	for i, pod := range pods.Items {
		nodeIDs[i] = pod.Name
	}
	return nodeIDs
}

// KillNode kills the given node
func KillNode(nodeID string) error {
	client := mustClient()
	pod := &corev1.Pod{}
	name := types.NamespacedName{
		Name:      nodeID,
		Namespace: GetNamespace(),
	}
	if err := client.Get(context.TODO(), name, pod); err != nil {
		return err
	}
	return client.Delete(context.TODO(), pod)
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
	return gnmi.New(ctx, dest)
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

// getConn gets a gRPC connection to the given address
func getConn(address string) (*grpc.ClientConn, error) {
	opts, err := handleCertArgs()
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
func GetAdminClient() (*grpc.ClientConn, proto.ConfigAdminServiceClient) {
	opts, err := handleCertArgs()
	if err != nil {
		fmt.Printf("Error loading cert %s", err)
	}
	conn, err := grpc.Dial(configAddress, opts...)
	if err != nil {
		panic(err)
	}
	return conn, proto.NewConfigAdminServiceClient(conn)
}

// NewAtomixClient returns an Atomix client from the environment
func NewAtomixClient(test string) (*atomixclient.Client, error) {
	opts := []atomixclient.ClientOption{
		atomixclient.WithNamespace(os.Getenv("ATOMIX_NAMESPACE")),
		atomixclient.WithApplication(test),
	}
	return atomixclient.NewClient(os.Getenv("ATOMIX_CONTROLLER"), opts...)
}

func mustClient() k8sclient.Client {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	kubeclient, err := k8sclient.New(config, k8sclient.Options{})
	if err != nil {
		panic(err)
	}
	return kubeclient
}
