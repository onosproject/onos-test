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

import corev1 "k8s.io/api/core/v1"

const (
	configType          = "config"
	configImage         = "onosproject/onos-config:latest"
	cmTestdeviceImageV1 = "onosproject/config-model-testdevice-1.0.0:latest"
	cmTestdeviceImageV2 = "onosproject/config-model-testdevice-2.0.0:latest"
	cmDevicesimImage    = "onosproject/config-model-devicesim-1.0.0:latest"
	cmStratumImage      = "onosproject/config-model-stratum-1.0.0:latest"
	configService       = "onos-config"
	configPort          = 5150
)

var configSecrets = map[string]string{
	"/certs/onf.cacrt":       caCert,
	"/certs/onos-config.crt": configCert,
	"/certs/onos-config.key": configKey,
}

const (
	volumeName         = "shared-data"
	volumePath         = "/usr/local/lib/shared"
	testDeviceNameV1   = "config-model-testdevice-1-0-0"
	testDeviceNameV2   = "config-model-testdevice-2-0-0"
	deviceSimName      = "config-model-devicesim-1-0-0"
	stratumName        = "config-model-stratum-1-0-0"
	modelPluginCommand = "/copylibandstay"
)

var testDeviceV1Args = []string{
	"testdevice.so.1.0.0",
	"/usr/local/lib/shared/testdevice.so.1.0.0",
	"stayrunning",
}
var testDeviceV2Args = []string{
	"testdevice.so.2.0.0",
	"/usr/local/lib/shared/testdevice.so.2.0.0",
	"stayrunning",
}

var deviceSimArgs = []string{
	"devicesim.so.1.0.0",
	"/usr/local/lib/shared/devicesim.so.1.0.0",
	"stayrunning",
}

var stratumArgs = []string{
	"stratum.so.1.0.0",
	"/usr/local/lib/shared/stratum.so.1.0.0",
	"stayrunning",
}

var configArgs = []string{
	"-caPath=/certs/onf.cacrt",
	"-keyPath=/certs/onos-config.key",
	"-certPath=/certs/onos-config.crt",
	"-modelPlugin=/usr/local/lib/shared/testdevice.so.1.0.0",
	"-modelPlugin=/usr/local/lib/shared/testdevice.so.2.0.0",
	"-modelPlugin=/usr/local/lib/shared/devicesim.so.1.0.0",
	"-modelPlugin=/usr/local/lib/shared/stratum.so.1.0.0",
}

func newConfig(cluster *Cluster) *Config {
	service := newService(cluster)
	ports := []Port{{Name: "grpc", Port: configPort}}
	serviceVolume := corev1.VolumeMount{Name: volumeName, MountPath: volumePath}
	service.SetArgs(configArgs...)
	service.SetSecrets(configSecrets)
	service.SetPorts(ports)
	service.SetLabels(getLabels(configType))
	service.SetImage(configImage)
	service.SetName(configService)
	service.SetVolume(serviceVolume)

	// Add model plugin sidecar containers
	var sideCars []*Sidecar

	// testdevice.so.1.0.0 model plugin
	container := newSidecar(cluster)
	container.SetName(testDeviceNameV1)
	container.SetImage(cmTestdeviceImageV1)
	container.SetArgs(testDeviceV1Args...)
	container.SetCommand(modelPluginCommand)
	volume := corev1.VolumeMount{Name: volumeName, MountPath: volumePath}
	container.SetVolume(volume)
	sideCars = append(sideCars, container)

	// testdevice.so.2.0.0 model plugin
	container = newSidecar(cluster)
	container.SetName(testDeviceNameV2)
	container.SetImage(cmTestdeviceImageV2)
	container.SetArgs(testDeviceV2Args...)
	container.SetCommand(modelPluginCommand)
	volume = corev1.VolumeMount{Name: volumeName, MountPath: volumePath}
	container.SetVolume(volume)
	sideCars = append(sideCars, container)

	// devicesim.so.1.0.0 model plugin
	container = newSidecar(cluster)
	container.SetName(deviceSimName)
	container.SetImage(cmDevicesimImage)
	container.SetArgs(deviceSimArgs...)
	container.SetCommand(modelPluginCommand)
	volume = corev1.VolumeMount{Name: volumeName, MountPath: volumePath}
	container.SetVolume(volume)
	sideCars = append(sideCars, container)

	// stratum.so.1.0.0 model plugin
	container = newSidecar(cluster)
	container.SetName(stratumName)
	container.SetImage(cmStratumImage)
	container.SetArgs(stratumArgs...)
	container.SetCommand(modelPluginCommand)
	volume = corev1.VolumeMount{Name: volumeName, MountPath: volumePath}
	container.SetVolume(volume)
	sideCars = append(sideCars, container)

	service.SetSidecars(sideCars)

	return &Config{
		Service: service,
	}
}

// Config provides methods for managing the onos-config service
type Config struct {
	*Service
}
