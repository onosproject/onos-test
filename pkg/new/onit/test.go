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

package onit

import (
	"github.com/onosproject/onos-test/pkg/new/kubetest"
	"github.com/onosproject/onos-test/pkg/new/onit/env"
	"github.com/onosproject/onos-test/pkg/new/onit/setup"
)

// Tests is the base type for ONIT test suites
type Tests struct {
	*kubetest.Tests
}

// Setup returns the ONOS setup API
func (t *Tests) Setup() setup.Setup {
	return setup.New(t.API())
}

// Env returns the ONOS environment API
func (t *Tests) Env() env.Env {
	return env.New(t.API())
}
