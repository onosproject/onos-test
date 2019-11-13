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

package kubetest

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"os"
	"path/filepath"
	"time"
)

const configPath = "/config"
const configFile = "config.yaml"

// TestType is the type of a test
type TestType string

const (
	// TestTypeTest is a type indicating a test
	TestTypeTest TestType = "test"
	// TestTypeBenchmark is a type indicating a benchmark
	TestTypeBenchmark TestType = "benchmark"
)

// TestConfig is a test configuration
type TestConfig struct {
	TestID     string
	Type       TestType
	Image      string
	Suite      string
	Test       string
	Timeout    time.Duration
	PullPolicy corev1.PullPolicy
}

// loadTestConfig loads the test configuration
func loadTestConfig() (*TestConfig, error) {
	file, err := os.Open(filepath.Join(configPath, configFile))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config := &TestConfig{}
	err = yaml.Unmarshal(jsonBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
