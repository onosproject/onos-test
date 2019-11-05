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
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
	"github.com/onosproject/onos-test/pkg/new/onit/deploy"
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

var _ Apps = &clusterApps{}

// clusterApps is an implementation of the Apps interface
type clusterApps struct {
	deployment deploy.Deployment
	apps       *cluster.Apps
}

func (e *clusterApps) List() []App {
	clusterApps := e.apps.List()
	apps := make([]App, len(clusterApps))
	for i, app := range clusterApps {
		apps[i] = e.Get(app.Name())
	}
	return apps
}

func (e *clusterApps) Get(name string) App {
	app := e.apps.Get(name)
	return &clusterApp{
		clusterService: &clusterService{
			service: app.Service,
		},
		app: app,
	}
}

func (e *clusterApps) Add(name string) deploy.App {
	return e.deployment.App(name)
}
