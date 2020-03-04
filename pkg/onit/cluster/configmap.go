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

func newConfigMap(cluster *Cluster) *ConfigMap {
	return &ConfigMap{}
}

// ConfigMap providing methods for adding configMaps to k8s cluster
type ConfigMap struct {
	name      string
	dataValue string
	dataKey   string
}

func (c *ConfigMap) DataKey() string {
	return c.dataKey
}

// SetDataKey sets data key
func (c *ConfigMap) SetDataKey(dataKey string) {
	c.dataKey = dataKey
}

// SetName sets a config map name
func (c *ConfigMap) SetName(name string) {
	c.name = name
}

// Name returns config map name
func (c *ConfigMap) Name() string {
	return c.name
}

// SetData sets config map data
func (c *ConfigMap) SetDataValue(data string) {
	c.dataValue = data
}

// Data returns configMap data
func (c *ConfigMap) DataValue() string {
	return c.dataValue
}
