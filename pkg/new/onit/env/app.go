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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// App provides the environment for an app
type App interface {
	Service
}

var _ App = &app{}

// app is an implementation of the App interface
type app struct {
	*service
}

func (e *app) Remove() {
	if err := e.teardownApp(); err != nil {
		panic(err)
	}
}

// teardownApp tears down a app by name
func (e *app) teardownApp() error {
	var err error
	if e := e.deleteAppDeployment(); e != nil {
		err = e
	}
	if e := e.deleteAppService(); e != nil {
		err = e
	}
	if e := e.deleteAppConfigMap(); e != nil {
		err = e
	}
	return err
}

// deleteAppConfigMap deletes an app ConfigMap by name
func (e *app) deleteAppConfigMap() error {
	return e.kubeClient.CoreV1().ConfigMaps(e.namespace).Delete(e.name, &metav1.DeleteOptions{})
}

// deleteAppPod deletes an app Pod by name
func (e *app) deleteAppDeployment() error {
	return e.kubeClient.AppsV1().Deployments(e.namespace).Delete(e.name, &metav1.DeleteOptions{})
}

// deleteAppService deletes an app Service by name
func (e *app) deleteAppService() error {
	return e.kubeClient.CoreV1().Services(e.namespace).Delete(e.name, &metav1.DeleteOptions{})
}
