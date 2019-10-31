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

package main

import (
	"github.com/onosproject/onos-test/pkg2/kubetest"
	_ "github.com/onosproject/onos-test/test/api"
	_ "github.com/onosproject/onos-test/test/atomix"
	_ "github.com/onosproject/onos-test/test/config"
	_ "github.com/onosproject/onos-test/test/integration"
	_ "github.com/onosproject/onos-test/test/topo"
	"github.com/onosproject/onos-test/test2"
)

func main() {
	kubetest.Register("suite-one", &test2.SuiteOne{})
	kubetest.Register("suite-two", &test2.SuiteTwo{})
	kubetest.Register("suite-three", &test2.SuiteThree{})
	kubetest.Main()
}
