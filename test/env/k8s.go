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
	"bytes"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	executil "k8s.io/client-go/util/exec"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
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

// GetCLINodes returns a list of onos-topo nodes
func GetCLINodes() []string {
	return getNodes(map[string]string{
		"app":  "onos",
		"type": "cli",
	})
}

// getNodes returns a list of nodes with the given labels
func getNodes(podLabels map[string]string) []string {
	kube := mustKubeClient()
	pods := &corev1.PodList{}
	options := &client.ListOptions{
		Namespace:     GetNamespace(),
		LabelSelector: labels.SelectorFromSet(podLabels),
	}
	err := kube.List(context.TODO(), options, pods)
	if err != nil {
		panic(err)
	}

	nodeIDs := make([]string, len(pods.Items))
	for i, pod := range pods.Items {
		nodeIDs[i] = pod.Name
	}
	return nodeIDs
}

// ExecuteCommand executes the given command on the given node (pod)
func ExecuteCommand(node string, command ...string) ([]string, int) {
	clientset := mustKubeClientset()
	pod, err := clientset.CoreV1().Pods(GetNamespace()).Get(node, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	container := pod.Spec.Containers[0]

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(node).
		Namespace(GetNamespace()).
		SubResource("exec").
		Param("container", container.Name)
	req.VersionedParams(&corev1.PodExecOptions{
		Container: container.Name,
		Command:   command,
		Stdout:    true,
		Stderr:    true,
		Stdin:     false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(mustKubeConfig(), "POST", req.URL())
	if err != nil {
		if execErr, ok := err.(executil.ExitError); ok && execErr.Exited() {
			return []string{}, execErr.ExitStatus()
		}
		panic(err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})

	if err != nil {
		if execErr, ok := err.(executil.ExitError); ok && execErr.Exited() {
			return []string{}, execErr.ExitStatus()
		}
		panic(err)
	}

	return strings.Split(strings.Trim(stdout.String(), "\n"), "\n"), 0
}

// KillNode kills the given node
func KillNode(nodeID string) error {
	client := mustKubeClient()
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

func mustKubeConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	return config
}

func mustKubeClientset() *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(mustKubeConfig())
	if err != nil {
		panic(err)
	}
	return clientset
}

func mustKubeClient() client.Client {
	kubeclient, err := client.New(mustKubeConfig(), client.Options{})
	if err != nil {
		panic(err)
	}
	return kubeclient
}
