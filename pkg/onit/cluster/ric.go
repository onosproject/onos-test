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
	ricType    = "ric"
	ricImage   = "onosproject/onos-ric:latest"
	ricService = "onos-ric"
	ricPort    = 5150
)

var ricSecrets = map[string]string{
	"/certs/onf.cacrt":    caCert,
	"/certs/onos-ric.crt": ricCert,
	"/certs/onos-ric.key": ricKey,
}

var ricArgs = []string{
	"-caPath=/certs/onf.cacrt",
	"-keyPath=/certs/onos-ric.key",
	"-certPath=/certs/onos-ric.crt",
}

// Enabled indicates whether the Ric is enabled
func (c *RIC) Enabled() bool {
	return GetArg(c.name, "enabled").Bool(c.enabled)
}

// SetEnabled sets whether the Ric is enabled
func (c *RIC) SetEnabled(enabled bool) {
	c.enabled = enabled
}

func newRIC(cluster *Cluster) *RIC {
	service := newService(cluster)
	ports := []Port{{Name: "grpc", Port: ricPort}}
	service.SetArgs(ricArgs...)
	service.SetSecrets(ricSecrets)
	service.SetPorts(ports)
	service.SetLabels(getLabels(ricType))
	service.SetImage(ricImage)
	service.SetName(ricService)

	return &RIC{
		Service: service,
	}
}

// RIC provides methods for managing the onos-ric service
type RIC struct {
	*Service
	enabled bool
}
