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
	"github.com/onosproject/onos-topo/pkg/northbound/admin"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
)

// TopoEnv provides the topo environment
type TopoEnv interface {
	ServiceEnv

	// NewAdminServiceClient returns the topo AdminService client
	NewAdminServiceClient() (admin.TopoAdminServiceClient, error)

	// NewDeviceServiceClient returns a topo device service client
	NewDeviceServiceClient() (device.DeviceServiceClient, error)
}

var _ TopoEnv = &clusterTopoEnv{}

// clusterTopoEnv is an implementation of the Topo interface
type clusterTopoEnv struct {
	*clusterServiceEnv
}

func (e *clusterTopoEnv) NewAdminServiceClient() (admin.TopoAdminServiceClient, error) {
	conn, err := e.Connect()
	if err != nil {
		return nil, err
	}
	return admin.NewTopoAdminServiceClient(conn), nil
}

func (e *clusterTopoEnv) NewDeviceServiceClient() (device.DeviceServiceClient, error) {
	conn, err := e.Connect()
	if err != nil {
		return nil, err
	}
	return device.NewDeviceServiceClient(conn), err
}
