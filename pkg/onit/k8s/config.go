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

package k8s

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	corev1 "k8s.io/api/core/v1"
)

var (
	_, path, _, _     = runtime.Caller(0)
	certsPath         = filepath.Join(filepath.Dir(filepath.Dir(path)), "../certs")
	configsPath       = filepath.Join(filepath.Dir(filepath.Dir(path)), "../../configs")
	deviceConfigsPath = filepath.Join(filepath.Join(filepath.Dir(filepath.Dir(path)), "../../configs"), "device")
	storeConfigsPath  = filepath.Join(filepath.Join(filepath.Dir(filepath.Dir(path)), "../../configs"), "store")
)

// ClusterConfig provides the configuration for the Kubernetes test cluster
type ClusterConfig struct {
	Registry      string            `yaml:"registry" mapstructure:"registry"`
	Preset        string            `yaml:"preset" mapstructure:"preset"`
	ImageTags     map[string]string `yaml:"image-tags" mapstructure:"image-tags"`
	PullPolicy    corev1.PullPolicy `yaml:"pull-policy" mapstructure:"pull-policy"`
	ConfigNodes   int               `yaml:"configNodes" mapstructure:"topoNodes"`
	TopoNodes     int               `yaml:"topoNodes" mapstructure:"topoNodes"`
	Partitions    int               `yaml:"partitions" mapstructure:"partitions"`
	PartitionSize int               `yaml:"partitionSize" mapstructure:"partitionSize"`
}

// load loads the preset configuration for the cluster
func (c *ClusterConfig) load() (map[string]interface{}, error) {
	file, err := os.Open(filepath.Join(storeConfigsPath, c.Preset+".json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var jsonObj map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonObj)
	if err != nil {
		return nil, err
	}
	return jsonObj, nil
}

// SimulatorConfig provides the configuration for a device simulator
type SimulatorConfig struct {
	Config string `yaml:"config" mapstructure:"config"`
}

// AppConfig provides the configuration for an app
type AppConfig struct {
	Image      string
	PullPolicy corev1.PullPolicy
}

// NetworkConfig provides the configuration for a stratum network
type NetworkConfig struct {
	Config         string `yaml:"config" mapstructure:"config"`
	MininetOptions []string
	NumDevices     int
	TopoType       TopoType
}

// load loads the simulator configuration
func (c *SimulatorConfig) load() (map[string]interface{}, error) {
	file, err := os.Open(filepath.Join(deviceConfigsPath, c.Config+".json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var jsonObj map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonObj)
	return jsonObj, err
}
