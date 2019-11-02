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

// NetworksEnv provides the networks environment
type NetworksEnv interface {
	Network(name string) NetworkEnv
	Add() setup.NetworkSetup
}

var _ NetworksEnv = &networksEnv{}

// networksEnv is an implementation of the NetworksEnv interface
type networksEnv struct {
	*testEnv
}

func (e *networksEnv) Network(name string) NetworkEnv {
	return &networkEnv{
		serviceEnv: &serviceEnv{
			testEnv: e.testEnv,
			name:    name,
		},
	}
}

func (e *networksEnv) Add() setup.NetworkSetup {
	return &networkSetup{
		testEnv: e.testEnv,
	}
}
