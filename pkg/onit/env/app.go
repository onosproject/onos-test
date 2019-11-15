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
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	corev1 "k8s.io/api/core/v1"
)

// AppSetup is an interface for setting up an application
type AppSetup interface {
	// Name sets the application name
	Name(name string) AppSetup

	// Nodes sets the number of application nodes
	Nodes(nodes int) AppSetup

	// Image sets the image to deploy
	Image(image string) AppSetup

	// PullPolicy sets the image pull policy
	PullPolicy(pullPolicy corev1.PullPolicy) AppSetup

	// Add adds the application to the cluster
	Add() (AppEnv, error)

	// AddOrDie adds the application and panics if the deployment fails
	AddOrDie() AppEnv
}

// clusterAppSetup is an implementation of the AppSetup interface
type clusterAppSetup struct {
	app *cluster.App
}

func (s *clusterAppSetup) Name(name string) AppSetup {
	s.app.SetName(name)
	return s
}

func (s *clusterAppSetup) Nodes(nodes int) AppSetup {
	s.app.SetReplicas(nodes)
	return s
}

func (s *clusterAppSetup) Image(image string) AppSetup {
	s.app.SetImage(image)
	return s
}

func (s *clusterAppSetup) PullPolicy(pullPolicy corev1.PullPolicy) AppSetup {
	s.app.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterAppSetup) Add() (AppEnv, error) {
	if err := s.app.Setup(); err != nil {
		return nil, err
	}
	return &clusterAppEnv{
		clusterServiceEnv: &clusterServiceEnv{
			service: s.app.Service,
		},
		app: s.app,
	}, nil
}

func (s *clusterAppSetup) AddOrDie() AppEnv {
	app, err := s.Add()
	if err != nil {
		panic(err)
	}
	return app
}

// AppEnv provides the environment for an app
type AppEnv interface {
	ServiceEnv

	// Remove removes the application
	Remove() error

	// RemoveOrDie removes the application and panics if the remove fails
	RemoveOrDie()
}

var _ AppEnv = &clusterAppEnv{}

// clusterAppEnv is an implementation of the App interface
type clusterAppEnv struct {
	*clusterServiceEnv
	app *cluster.App
}

func (e *clusterAppEnv) Remove() error {
	return e.app.TearDown()
}

func (e *clusterAppEnv) RemoveOrDie() {
	if err := e.Remove(); err != nil {
		panic(err)
	}
}
