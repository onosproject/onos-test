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
	atomixcontroller "github.com/atomix/atomix-k8s-controller/pkg/client/clientset/versioned"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// client is the base for all Kubernetes cluster objects
type client struct {
	namespace        string
	kubeClient       *kubernetes.Clientset
	atomixClient     *atomixcontroller.Clientset
	extensionsClient *apiextension.Clientset
}

func (c *client) listPods(serviceType serviceType) []string {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{typeLabel: serviceType.name()}}
	pods, err := c.kubeClient.CoreV1().Pods(c.namespace).List(metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String()})
	if err != nil {
		panic(err)
	}

	names := make([]string, len(pods.Items))
	for i, dep := range pods.Items {
		names[i] = dep.Name
	}
	return names
}

func (c *client) listServices(serviceType serviceType) []string {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{typeLabel: serviceType.name()}}
	services, err := c.kubeClient.CoreV1().Services(c.namespace).List(metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String()})
	if err != nil {
		panic(err)
	}

	names := make([]string, len(services.Items))
	for i, dep := range services.Items {
		names[i] = dep.Name
	}
	return names
}

func (c *client) listDeployments(serviceType serviceType) []string {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{typeLabel: serviceType.name()}}
	deps, err := c.kubeClient.AppsV1().Deployments(c.namespace).List(metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String()})
	if err != nil {
		panic(err)
	}

	names := make([]string, len(deps.Items))
	for i, dep := range deps.Items {
		names[i] = dep.Name
	}
	return names
}
