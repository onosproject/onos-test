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
	"context"
	"github.com/onosproject/onos-config/api/admin"
	"github.com/onosproject/onos-config/api/diags"
	"github.com/openconfig/gnmi/client"
	gnmi "github.com/openconfig/gnmi/client/gnmi"
	"time"
)

// ConfigEnv provides the config environment
type ConfigEnv interface {
	ServiceEnv

	// Destination returns the gNMI client destination
	Destination() client.Destination

	// NewAdminServiceClient returns the config AdminService client
	NewAdminServiceClient() (admin.ConfigAdminServiceClient, error)

	// NewChangeServiceClient returns the config AdminService client
	NewChangeServiceClient() (diags.ChangeServiceClient, error)

	// NewGNMIClient returns the gNMI client
	NewGNMIClient() (*gnmi.Client, error)
}

var _ ConfigEnv = &clusterConfigEnv{}

// clusterConfigEnv is an implementation of the Config interface
type clusterConfigEnv struct {
	*clusterServiceEnv
}

func (e *clusterConfigEnv) Destination() client.Destination {
	return client.Destination{
		Addrs:   []string{e.Address()},
		Target:  "gnmi",
		TLS:     e.Credentials(),
		Timeout: 10 * time.Second,
	}
}

func (e *clusterConfigEnv) NewAdminServiceClient() (admin.ConfigAdminServiceClient, error) {
	conn, err := e.Connect()
	if err != nil {
		return nil, err
	}
	return admin.NewConfigAdminServiceClient(conn), nil
}

func (e *clusterConfigEnv) NewGNMIClient() (*gnmi.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	gnmiClient, err := gnmi.New(ctx, e.Destination())
	if err != nil {
		return nil, err
	}
	return gnmiClient.(*gnmi.Client), nil
}

func (e *clusterConfigEnv) NewChangeServiceClient() (diags.ChangeServiceClient, error) {
	conn, err := e.Connect()
	if err != nil {
		return nil, err
	}
	return diags.NewChangeServiceClient(conn), nil
}
