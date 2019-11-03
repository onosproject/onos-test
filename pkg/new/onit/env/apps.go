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

import "github.com/onosproject/onos-test/pkg/new/onit/setup"

// Apps provides the environment for applications
type Apps interface {
	// Apps returns a list of all apps in the environment
	Apps() []App

	// Get returns the environment for an app by name
	Get(name string) App

	// Add adds an app to the environment
	Add(name string) setup.AppSetup
}

var _ Apps = &apps{}

// apps is an implementation of the Apps interface
type apps struct {
	*testEnv
}

func (e *apps) Apps() []App {
	panic("implement me")
}

func (e *apps) Get(name string) App {
	return &app{
		service: &service{
			testEnv: e.testEnv,
			name:    name,
		},
	}
}

func (e *apps) Add(name string) setup.AppSetup {
	return &appSetup{
		serviceTypeSetup: &serviceTypeSetup{
			serviceSetup: &serviceSetup{
				testEnv: e.testEnv,
			},
		},
		name: name,
	}
}
