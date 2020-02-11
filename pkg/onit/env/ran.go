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
	"github.com/onosproject/onos-ran/api/nb"
	"google.golang.org/grpc"
)

// RanEnv provides the topo environment
type RanEnv interface {
	ServiceEnv

	// NewRanC1ServiceClient returns a RAN C1 service client
	NewRanC1ServiceClient() (nb.C1InterfaceServiceClient, error)
}

var _ RanEnv = &clusterRanEnv{}

// clusterRanEnv is an implementation of the RAN interface
type clusterRanEnv struct {
	*clusterServiceEnv
}

func (e *clusterRanEnv) NewRanC1ServiceClient() (nb.C1InterfaceServiceClient, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, connErr := grpc.Dial(env.Ran().Address(), opts...)
	if connErr != nil {
		return nil, connErr
	}

	return nb.NewC1InterfaceServiceClient(conn), nil
}
