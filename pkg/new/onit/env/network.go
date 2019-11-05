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
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// Network provides the environment for a network node
type Network interface {
	Service

	// Devices returns a list of devices in the network
	Devices() []Node
}

var _ Network = &network{}

// network is an implementation of the Network interface
type network struct {
	*service
}

func (e *network) Devices() []Node {
	services, err := e.kubeClient.CoreV1().Services(e.namespace).List(metav1.ListOptions{
		LabelSelector: "type=network,network=" + e.name,
	})
	if err != nil {
		panic(err)
	}

	devices := make([]Node, len(services.Items))
	for i, service := range services.Items {
		devices[i] = &node{
			testEnv: e.testEnv,
			name:    service.Name,
		}
	}
	return devices
}

func (e *network) Remove() {
	if err := e.teardownNetwork(); err != nil {
		panic(err)
	}
}

// teardownNetwork tears down a network by name
func (e *network) teardownNetwork() error {
	pods, err := e.kubeClient.CoreV1().Pods(e.namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("type=network,network=%s", e.name),
	})
	if err != nil {
		return err
	} else if len(pods.Items) == 0 {
		return fmt.Errorf("no resources matching '%s' found", e.name)
	}
	total := len(pods.Items)
	if e := e.deletePod(); e != nil {
		err = e
	}
	if e := e.deleteService(); e != nil {
		err = e
	}

	for total > 0 {
		time.Sleep(50 * time.Millisecond)
		pods, err = e.kubeClient.CoreV1().Pods(e.namespace).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf("type=network,network=%s", e.name),
		})
		if err != nil {
			return err
		}

		total = len(pods.Items)

	}
	return err
}

// deletePod deletes a network Pod by name
func (e *network) deletePod() error {
	return e.kubeClient.CoreV1().Pods(e.namespace).Delete(e.name, &metav1.DeleteOptions{})
}

// deleteService deletes all network Service by name
func (e *network) deleteService() error {
	label := "type=network,network=" + e.name
	serviceList, _ := e.kubeClient.CoreV1().Services(e.namespace).List(metav1.ListOptions{
		LabelSelector: label,
	})

	for _, svc := range serviceList.Items {
		err := e.kubeClient.CoreV1().Services(e.namespace).Delete(svc.Name, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
