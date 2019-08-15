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

package onit

import (
	"os"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"golang.org/x/crypto/ssh/terminal"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenShell opens a shell session to the given resource
func (c *ClusterController) OpenShell(resourceID string) error {
	pod, err := c.kubeclient.CoreV1().Pods(c.clusterID).Get(resourceID, metav1.GetOptions{})
	if err != nil {
		return err
	}

	container := pod.Spec.Containers[0]
	req := c.kubeclient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec").
		Param("container", container.Name)
	req.VersionedParams(&corev1.PodExecOptions{
		Container: container.Name,
		Command:   []string{"/bin/sh"},
		Stdout:    true,
		Stdin:     true,
		TTY:       true,
	}, scheme.ParameterCodec)

	config, err := getRestConfig()
	if err != nil {
		return err
	}

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := terminal.Restore(0, oldState)
		if err != nil {
			panic(err)
		}

	}()

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Tty:    true,
	})
	return err
}

// setupOnosCli sets up the onos-cli deployment
func (c *ClusterController) setupOnosCli() error {
	if err := c.createCliDeployment(); err != nil {
		return err
	}
	if err := c.awaitCliDeploymentReady(); err != nil {
		return err
	}
	return nil
}

// createCliDeployment creates an onos-cli deployment
func (c *ClusterController) createCliDeployment() error {
	nodes := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-cli",
			Namespace: c.clusterID,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nodes,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  "onos",
					"type": "cli",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  "onos",
						"type": "cli",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "onos-cli",
							Image:           c.imageName("onosproject/onos-cli", c.config.ImageTags["cli"]),
							ImagePullPolicy: c.config.PullPolicy,
							Stdin:           true,
						},
					},
				},
			},
		},
	}
	_, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Create(deployment)
	return err
}

// awaitCliDeploymentReady waits for the onos-cli pods to complete startup
func (c *ClusterController) awaitCliDeploymentReady() error {
	for {
		// Get the onos-cli deployment
		dep, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Get("onos-cli", metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Return once the all replicas in the deployment are ready
		if int(dep.Status.ReadyReplicas) == 1 {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// GetOnosTopoNodes returns a list of all onos-topo nodes running in the cluster
func (c *ClusterController) GetOnosCliNodes() ([]NodeInfo, error) {
	topoLabelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": "onos", "type": "cli"}}

	pods, err := c.kubeclient.CoreV1().Pods(c.clusterID).List(metav1.ListOptions{
		LabelSelector: labels.Set(topoLabelSelector.MatchLabels).String(),
	})
	if err != nil {
		return nil, err
	}

	onosCliNodes := make([]NodeInfo, len(pods.Items))
	for i, pod := range pods.Items {
		var status NodeStatus
		if pod.Status.Phase == corev1.PodRunning {
			status = NodeRunning
		} else if pod.Status.Phase == corev1.PodFailed {
			status = NodeFailed
		}
		onosCliNodes[i] = NodeInfo{
			ID:     pod.Name,
			Status: status,
			Type:   OnosCli,
		}
	}

	return onosCliNodes, nil
}
