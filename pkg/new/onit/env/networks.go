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

// Networks provides the networks environment
type Networks interface {
	// Networks returns a list of networks in the environment
	Networks() []Network

	// Get returns the environment for a network service by name
	Get(name string) Network

	// Add adds a new network to the environment
	Add(name string) setup.NetworkSetup
}

var _ Networks = &networks{}

// networks is an implementation of the Networks interface
type networks struct {
	*testEnv
}

func (e *networks) Networks() []Network {
	panic("implement me")
}

func (e *networks) Get(name string) Network {
	return &network{
		service: &service{
			testEnv: e.testEnv,
			name:    name,
		},
	}
}

func (e *networks) Add(name string) setup.NetworkSetup {
	return &networkSetup{
		serviceTypeSetup: &serviceTypeSetup{
			serviceSetup: &serviceSetup{
				testEnv: e.testEnv,
			},
		},
		name: name,
	}
}
