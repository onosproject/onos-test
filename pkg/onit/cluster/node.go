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
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	executil "k8s.io/client-go/util/exec"
)

func newNode(cluster *Cluster) *Node {
	return &Node{
		client:     cluster.client,
		cluster:    cluster,
		pullPolicy: corev1.PullIfNotPresent,
	}
}

// Node provides the environment for a single node
type Node struct {
	*client
	cluster    *Cluster
	name       string
	port       int
	image      string
	pullPolicy corev1.PullPolicy
}

// Name returns the node name
func (n *Node) Name() string {
	return GetArg(n.name, "service").String(n.name)
}

// SetName sets the node name
func (n *Node) SetName(name string) {
	n.name = name
}

// SetPort sets the node port
func (n *Node) SetPort(port int) {
	n.port = port
}

func (n *Node) Port() int {
	return n.port
}

// Address returns the service address
func (n *Node) Address() string {
	return fmt.Sprintf("%s:%d", n.name, n.port)
}

// Image returns the image configured for the node
func (n *Node) Image() string {
	return GetArg(n.name, "image").String(n.image)
}

// SetImage sets the node image
func (n *Node) SetImage(image string) {
	n.image = image
}

// PullPolicy returns the image pull policy configured for the node
func (n *Node) PullPolicy() corev1.PullPolicy {
	return corev1.PullPolicy(GetArg(n.name, "pullPolicy").String(string(n.pullPolicy)))
}

// SetPullPolicy sets the image pull policy for the node
func (n *Node) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	n.pullPolicy = pullPolicy
}

// Execute executes the given command on the node
func (n *Node) Execute(command ...string) ([]string, int, error) {
	pod, err := n.kubeClient.CoreV1().Pods(n.namespace).Get(n.Name(), metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	container := pod.Spec.Containers[0]

	fullCommand := append([]string{"/bin/bash", "-c"}, command...)
	req := n.kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(n.Name()).
		Namespace(n.namespace).
		SubResource("exec").
		Param("container", container.Name)
	req.VersionedParams(&corev1.PodExecOptions{
		Container: container.Name,
		Command:   fullCommand,
		Stdout:    true,
		Stderr:    true,
		Stdin:     false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(n.config, "POST", req.URL())
	if err != nil {
		if execErr, ok := err.(executil.ExitError); ok && execErr.Exited() {
			return []string{}, execErr.ExitStatus(), nil
		}
		return nil, 0, err
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})

	if err != nil {
		if execErr, ok := err.(executil.ExitError); ok && execErr.Exited() {
			return []string{}, execErr.ExitStatus(), nil
		}
		return nil, 0, err
	}

	return strings.Split(strings.Trim(stdout.String(), "\n"), "\n"), 0, nil
}

// Credentials returns the TLS credentials
func (n *Node) Credentials() (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}, nil
}

// Connect creates a gRPC client connection to the node
func (n *Node) Connect() (*grpc.ClientConn, error) {
	tlsConfig, err := n.Credentials()
	if err != nil {
		return nil, err
	}
	return grpc.Dial(n.Address(), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}

// AwaitReady waits for the node to become ready
func (n *Node) AwaitReady() error {
	for {
		ready, err := n.isReady()
		if err != nil {
			return err
		} else if ready {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// isReady returns a bool indicating whether the node is ready
func (n *Node) isReady() (bool, error) {
	pod, err := n.kubeClient.CoreV1().Pods(n.namespace).Get(n.name, metav1.GetOptions{})
	if err != nil {
		return false, err
	} else if pod == nil {
		return false, errors.New("node not found")
	}

	for _, status := range pod.Status.ContainerStatuses {
		if !status.Ready {
			return false, nil
		}
	}
	return true, nil
}

// Delete deletes the node
func (n *Node) Delete() error {
	return n.kubeClient.CoreV1().Pods(n.namespace).Delete(n.name, &metav1.DeleteOptions{})
}
