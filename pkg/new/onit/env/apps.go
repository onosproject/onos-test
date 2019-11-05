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
	"github.com/onosproject/onos-test/pkg/new/onit/deploy"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// Apps provides the environment for applications
type Apps interface {
	// List returns a list of all apps in the environment
	List() []App

	// Get returns the environment for an app by name
	Get(name string) App

	// Add adds an app to the environment
	Add(name string) deploy.App
}

var _ Apps = &apps{}

// apps is an implementation of the Apps interface
type apps struct {
	*testEnv
}

// GetApps returns a list of apps deployed in the cluster
func (e *apps) List() []App {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"type": "app"}}
	appList, err := e.kubeClient.AppsV1().Deployments(e.namespace).List(metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String()})
	if err != nil {
		panic(err)
	}

	apps := make([]App, len(appList.Items))
	for i, app := range appList.Items {
		apps[i] = e.Get(app.Name)
	}
	return apps
}

func (e *apps) Get(name string) App {
	return &app{
		service: newService(name, "app", e.testEnv),
	}
}

func (e *apps) Add(name string) deploy.App {
	return e.deployment.App(name)
}
