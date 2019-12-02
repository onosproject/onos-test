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

package test

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

// JobType is the type of a test
type JobType string

const (
	// TestTypeTest is a type indicating a test
	TestTypeTest JobType = "test"
	// TestTypeBenchmark is a type indicating a benchmark
	TestTypeBenchmark JobType = "benchmark"
)

// JobConfig is a job configuration
type JobConfig struct {
	JobID      string
	Type       JobType
	Image      string
	Env        map[string]string
	Timeout    time.Duration
	PullPolicy corev1.PullPolicy
	Teardown   bool
}

// Config is a configuration that provides a job config
type Config interface {
	// Job returns the job configuration
	Job() *JobConfig
}

// TestConfig is a test configuration
type TestConfig struct {
	*JobConfig
	Suite string
	Test  string
}

func (c *TestConfig) Job() *JobConfig {
	return c.JobConfig
}

// BenchmarkConfig is a benchmark configuration
type BenchmarkConfig struct {
	*JobConfig
	Suite       string
	Benchmark   string
	Clients     int
	Parallelism int
	Requests    int
	Args        map[string]string
}

func (c *BenchmarkConfig) Job() *JobConfig {
	return c.JobConfig
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

// loadBenchmarkConfig loads the test configuration
func loadBenchmarkConfig() (*BenchmarkConfig, error) {
	file, err := os.Open(filepath.Join(configPath, configFile))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config := &BenchmarkConfig{}
	err = yaml.Unmarshal(jsonBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
