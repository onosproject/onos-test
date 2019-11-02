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

// ServiceEnv is a base interface for service environments
type ServiceEnv interface {
	// Name is the name of the service
	Name() string

	// Nodes returns the service nodes
	Nodes() []NodeEnv

	// Node returns a specific node environment
	Node(name string) NodeEnv

	// Remove removes the service
	Remove()
}

var _ ServiceEnv = &serviceEnv{}

// serviceEnv is an implementation of the ServiceEnv interface
type serviceEnv struct {
	*testEnv
	name string
}

func (e *serviceEnv) Name() string {
	return e.name
}

func (e *serviceEnv) Nodes() []NodeEnv {
	panic("implement me")
}

func (e *serviceEnv) Node(name string) NodeEnv {
	return &nodeEnv{
		testEnv: e.testEnv,
		name:    name,
	}
}

func (e *serviceEnv) Remove() {
	panic("implement me")
}
