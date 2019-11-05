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

// Simulator provides the environment for a single simulator
type Simulator interface {
	Service
}

var _ Simulator = &simulator{}

// simulator is an implementation of the Simulator interface
type simulator struct {
	*service
}

func (s *simulator) Remove() {
	if err := s.teardownSimulator(); err != nil {
		panic(err)
	}
}

// teardownSimulator tears down a simulator by name
func (s *simulator) teardownSimulator() error {
	pods, err := s.kubeClient.CoreV1().Pods(s.namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("type=simulator,simulator=%s", s.name),
	})
	if err != nil {
		return err
	} else if len(pods.Items) == 0 {
		return fmt.Errorf("no resources matching '%s' found", s.name)
	}
	total := len(pods.Items)

	if e := s.deletePod(); e != nil {
		err = e
	}
	if e := s.deleteService(); e != nil {
		err = e
	}
	if e := s.deleteConfigMap(); e != nil {
		err = e
	}

	for total > 0 {
		time.Sleep(50 * time.Millisecond)
		pods, err = s.kubeClient.CoreV1().Pods(s.namespace).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf("type=simulator,simulator=%s", s.name),
		})
		if err != nil {
			return err
		}

		total = len(pods.Items)

	}
	return err
}

// deleteConfigMap deletes a simulator ConfigMap by name
func (s *simulator) deleteConfigMap() error {
	return s.kubeClient.CoreV1().ConfigMaps(s.namespace).Delete(s.name, &metav1.DeleteOptions{})
}

// deletePod deletes a simulator Pod by name
func (s *simulator) deletePod() error {
	return s.kubeClient.CoreV1().Pods(s.namespace).Delete(s.name, &metav1.DeleteOptions{})
}

// deleteService deletes a simulator Service by name
func (s *simulator) deleteService() error {
	return s.kubeClient.CoreV1().Services(s.namespace).Delete(s.name, &metav1.DeleteOptions{})
}
