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
	guiType    = "gui"
	guiImage   = "onosproject/onos-gui:latest"
	guiService = "onos-gui"
	guiPort    = 80
)

// Enabled indicates whether the Gui is enabled
func (c *Gui) Enabled() bool {
	return GetArg(c.name, "enabled").Bool(c.enabled)
}

// SetEnabled sets whether the Gui is enabled
func (c *Gui) SetEnabled(enabled bool) {
	c.enabled = enabled
}

func newGui(cluster *Cluster) *Gui {
	return &Gui{
		Service: newService(cluster, guiService, []Port{{Name: "grpc", Port: guiPort}}, getLabels(guiType), guiImage, nil, nil, nil),
	}
}

// Gui provides methods for managing the onos-gui service
type Gui struct {
	*Service
	enabled bool
}
