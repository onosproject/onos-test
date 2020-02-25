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

package cluster

const (
	ranType    = "ran"
	ranImage   = "onosproject/onos-ric:latest"
	ranService = "onos-ric"
	ranPort    = 5150
)

var ranSecrets = map[string]string{
	"/certs/onf.cacrt":    caCert,
	"/certs/onos-ric.crt": ranCert,
	"/certs/onos-ric.key": ranKey,
}

var ranArgs = []string{
	"-caPath=/certs/onf.cacrt",
	"-keyPath=/certs/onos-ric.key",
	"-certPath=/certs/onos-ric.crt",
	"-simulator=ran-simulator:5150",
}

// Enabled indicates whether the Ran is enabled
func (c *RAN) Enabled() bool {
	return GetArg(c.name, "enabled").Bool(c.enabled)
}

// SetEnabled sets whether the Ran is enabled
func (c *RAN) SetEnabled(enabled bool) {
	c.enabled = enabled
}

func newRAN(cluster *Cluster) *RAN {
	service := newService(cluster)
	ports := []Port{{Name: "grpc", Port: ranPort}}
	service.SetArgs(ranArgs...)
	service.SetSecrets(ranSecrets)
	service.SetPorts(ports)
	service.SetLabels(getLabels(ranType))
	service.SetImage(ranImage)
	service.SetName(ranService)

	return &RAN{
		Service: service,
	}
}

// RAN provides methods for managing the onos-ric service
type RAN struct {
	*Service
	enabled bool
}
