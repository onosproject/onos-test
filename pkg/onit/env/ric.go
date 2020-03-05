// Copyright 2020-present Open Networking Foundation.
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
	"github.com/onosproject/onos-ric/api/nb"
)

// RICEnv provides the RIC environment
type RICEnv interface {
	ServiceEnv

	// NewRICC1ServiceClient returns a RIC C1 service client
	NewRICC1ServiceClient() (nb.C1InterfaceServiceClient, error)
}

var _ RICEnv = &clusterRICEnv{}

// clusterRICEnv is an implementation of the RAN interface
type clusterRICEnv struct {
	*clusterServiceEnv
}

func (e *clusterRICEnv) NewRICC1ServiceClient() (nb.C1InterfaceServiceClient, error) {
	conn, connErr := e.Connect()
	if connErr != nil {
		return nil, connErr
	}

	return nb.NewC1InterfaceServiceClient(conn), nil
}
