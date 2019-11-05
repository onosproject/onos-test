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

package deploy

import (
	"github.com/onosproject/onos-test/pkg/new/onit/cluster"
)

// App is an interface for setting up an application
type App interface {
	Deploy
	ServiceType

	// Nodes sets the number of application nodes
	Nodes(nodes int) App
}

// clusterApp is an implementation of the App interface
type clusterApp struct {
	*clusterServiceType
	app *cluster.App
}

func (s *clusterApp) Nodes(nodes int) App {
	s.app.SetNodes(nodes)
	return s
}

func (s *clusterApp) Deploy() error {
	return s.app.Add()
}

func (s *clusterApp) DeployOrDie() {
	if err := s.Deploy(); err != nil {
		panic(err)
	}
}
