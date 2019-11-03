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

// Node provides the environment for a single node
type Node interface {
	// Name returns the name of the node
	Name() string

	// Kill kills the node
	Kill()
}

// node is an implementation of the Node interface
type node struct {
	*testEnv
	name string
}

func (e *node) Name() string {
	return e.name
}

func (e *node) Kill() {
	panic("implement me")
}
