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

import (
	"gopkg.in/yaml.v2"
)

// NewTracingConfig creates an instance of tracing configuration
func NewTracingConfig(data []byte) *TracingConfiguration {
	return &TracingConfiguration{
		data: data,
	}
}

// NewSimulatorConfig creates an instance of simulator configuration
func NewSimulatorConfig(data []byte) *SimulatorConfiguration {
	return &SimulatorConfiguration{
		data: data,
	}
}

// SimulatorConfiguration config data structure for simulator configuration info
type SimulatorConfiguration struct {
	config SimulatorYamlConfig
	data   []byte
}

// TracingConfiguration config data structure for tracing config file
type TracingConfiguration struct {
	config TracingYamlConfig
	data   []byte
}

// SimulatorYamlConfig yaml data structure for the simulator  configuration file
type SimulatorYamlConfig struct {
	Configuration string `yaml:"configuration"`
}

// TracingYamlConfig yaml data structure for the tracing configuration file
type TracingYamlConfig struct {
	Tracing struct {
		Logging struct {
			Loggers []struct {
				Encoding string `yaml:"encoding"`
				Level    string `yaml:"level"`
				Name     string `yaml:"name"`
				Sink     string `yaml:"sink"`
			} `yaml:"loggers"`
			Sinks []struct {
				Key   string `yaml:"key"`
				Name  string `yaml:"name"`
				Topic string `yaml:"topic"`
				Type  string `yaml:"type"`
				URI   string `yaml:"uri"`
			} `yaml:"sinks"`
		} `yaml:"logging"`
	} `yaml:"tracing"`
}

// GetConfig return the yaml config data structure
func (c *TracingConfiguration) GetConfig() TracingYamlConfig {
	return c.config
}

// Parse parse a yaml config file
func (c *TracingConfiguration) Parse() error {
	err := yaml.Unmarshal(c.data, &c.config)
	return err
}

// GetConfig return the yaml config data structure
func (c *SimulatorConfiguration) GetConfig() SimulatorYamlConfig {
	return c.config
}

// Parse parse a yaml config file
func (c *SimulatorConfiguration) Parse() error {
	err := yaml.Unmarshal(c.data, &c.config)
	return err
}
