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

package cluster

import "github.com/onosproject/onos-test/pkg/new/util/random"

func newApps(client *client) *Apps {
	return &Apps{
		client: client,
	}
}

// Apps provides methods for adding and modifying applications
type Apps struct {
	*client
}

// New returns a new app
func (s *Apps) New() *App {
	return newApp(random.NewPetName(2), s.client)
}

// Get gets an app by name
func (s *Apps) Get(name string) *App {
	return newApp(name, s.client)
}

// List lists the networks in the cluster
func (s *Apps) List() []*App {
	names := s.listDeployments(appType)
	apps := make([]*App, len(names))
	for i, name := range names {
		apps[i] = s.Get(name)
	}
	return apps
}
