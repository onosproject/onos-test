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

const (
	envoyType    = "envoy"
	envoyService = "onos-envoy"
	envoyImage   = "envoyproxy/envoy-alpine:latest"
	envoyPort    = 8080
)

var envoyCommand = []string{
	"/usr/local/bin/envoy",
	"-c",
	"/etc/envoy-proxy/config/envoy-config.yaml",
}

var envoySecrets = map[string]string{
	"/certs/onf.cacrt":  caCert,
	"/certs/client.crt": clientCert,
	"/certs/client.key": clientKey,
}

// Enabled indicates whether the Gui is enabled
func (c *Envoy) Enabled() bool {
	return GetArg(c.name, "enabled").Bool(c.enabled)
}

// SetEnabled sets whether the Envoy is enabled
func (c *Envoy) SetEnabled(enabled bool) {
	c.enabled = enabled
}

func newEnvoy(cluster *Cluster) *Envoy {
	service := newService(cluster, envoyService, []Port{{Name: "envoy", Port: envoyPort}}, getLabels(envoyType), envoyImage, envoySecrets, nil)
	service.SetCommand(envoyCommand...)
	return &Envoy{
		Service: service,
	}
}

// Envoy provides methods for managing the envoy service
type Envoy struct {
	*Service
	enabled bool
}
