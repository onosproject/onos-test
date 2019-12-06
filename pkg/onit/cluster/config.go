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
	configType    = "config"
	configImage   = "onosproject/onos-config:latest"
	configService = "onos-config"
	configPort    = 5150
)

var configSecrets = map[string]string{
	"/certs/onf.cacrt":       caCert,
	"/certs/onos-config.crt": configCert,
	"/certs/onos-config.key": configKey,
}

var configArgs = []string{
	"-caPath=/certs/onf.cacrt",
	"-keyPath=/certs/onos-config.key",
	"-certPath=/certs/onos-config.crt",
	"-modelPlugin=/usr/local/lib/testdevice.so.1.0.0",
	"-modelPlugin=/usr/local/lib/testdevice.so.2.0.0",
	"-modelPlugin=/usr/local/lib/devicesim.so.1.0.0",
	"-modelPlugin=/usr/local/lib/stratum.so.1.0.0",
}

func newConfig(cluster *Cluster) *Config {
	return &Config{
		Service: newService(cluster, configService, []Port{{Name: "grpc", Port: configPort}}, getLabels(configType), configImage, configSecrets, configArgs, nil, nil),
	}
}

// Config provides methods for managing the onos-config service
type Config struct {
	*Service
}
