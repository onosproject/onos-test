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

package k8s

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"

	atomixk8s "github.com/atomix/atomix-k8s-controller/pkg/client/clientset/versioned"
	"github.com/onosproject/onos-test/pkg/onit/console"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// ClusterController manages a single cluster in Kubernetes
type ClusterController struct {
	ClusterID        string
	Restconfig       *rest.Config
	Kubeclient       *kubernetes.Clientset
	Atomixclient     *atomixk8s.Clientset
	Extensionsclient *apiextension.Clientset
	Config           *ClusterConfig
	Status           *console.StatusWriter
}

// Setup sets up a test cluster with the given configuration
func (c *ClusterController) Setup() console.ErrorStatus {
	c.Status.Start("Setting up RBAC")
	if err := c.setupRBAC(); err != nil {
		return c.Status.Fail(err)
	}
	c.Status.Succeed()
	c.Status.Start("Setting up Atomix controller")
	if err := c.setupAtomixController(); err != nil {
		return c.Status.Fail(err)
	}
	c.Status.Succeed()
	c.Status.Start("Starting Raft partitions")
	if err := c.setupPartitions(); err != nil {
		return c.Status.Fail(err)
	}
	c.Status.Succeed()
	c.Status.Start("Creating secret for onos nodes")
	if err := c.createOnosSecret(); err != nil {
		return c.Status.Fail(err)
	}
	c.Status.Succeed()
	c.Status.Start("Bootstrapping onos-topo cluster")
	if err := c.setupOnosTopo(); err != nil {
		return c.Status.Fail(err)
	}
	c.Status.Succeed()
	c.Status.Start("Bootstrapping onos-config cluster")
	if err := c.setupOnosConfig(); err != nil {
		return c.Status.Fail(err)
	}
	return c.Status.Succeed()
}

func (c *ClusterController) Teardown() console.ErrorStatus {
	//TODO needs to be implemented
	return c.Status.Succeed()
}

// GetResources returns a list of resource IDs matching the given resource name
func (c *ClusterController) GetResources(name string) ([]string, error) {
	pod, err := c.Kubeclient.CoreV1().Pods(c.ClusterID).Get(name, metav1.GetOptions{})
	if err == nil {
		return []string{pod.Name}, nil
	} else if !k8serrors.IsNotFound(err) {
		return nil, err
	}

	pods, err := c.Kubeclient.CoreV1().Pods(c.ClusterID).List(metav1.ListOptions{
		LabelSelector: "resource=" + name,
	})
	if err != nil {
		return nil, err
	} else if len(pods.Items) == 0 {
		return nil, errors.New("unknown test resource " + name)
	}

	resources := make([]string, len(pods.Items))
	for i, pod := range pods.Items {
		resources[i] = pod.Name
	}
	return resources, nil
}

// PortForward forwards a local port to the given remote port on the given resource
func (c *ClusterController) PortForward(resourceID string, localPort int, remotePort int) error {
	pod, err := c.Kubeclient.CoreV1().Pods(c.ClusterID).Get(resourceID, metav1.GetOptions{})
	if err != nil {
		return err
	}

	req := c.Kubeclient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("portforward")

	roundTripper, upgradeRoundTripper, err := spdy.RoundTripperFor(c.Restconfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgradeRoundTripper, &http.Client{Transport: roundTripper}, http.MethodPost, req.URL())

	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", localPort, remotePort)}, stopChan, readyChan, out, errOut)
	if err != nil {
		return err
	}

	go func() {
		for range readyChan { // Kubernetes will close this channel when it has something to tell us.
		}
		if len(errOut.String()) != 0 {
			fmt.Println(errOut.String())
			os.Exit(1)
		} else if len(out.String()) != 0 {
			fmt.Println(out.String())
		}
	}()

	return forwarder.ForwardPorts()
}
