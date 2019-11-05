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
)

// App provides the environment for an app
type App interface {
	Service

	// Remove removes the application
	Remove() error

	// RemoveOrDie removes the application and panics if the remove fails
	RemoveOrDie()
}

var _ App = &clusterApp{}

// clusterApp is an implementation of the App interface
type clusterApp struct {
	*clusterService
	app *cluster.App
}

func (e *clusterApp) Remove() error {
	return e.app.Remove()
}

func (e *clusterApp) RemoveOrDie() {
	if err := e.Remove(); err != nil {
		panic(err)
	}
}
