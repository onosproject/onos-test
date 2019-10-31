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
	TestTypeTest      TestType = "test"
	TestTypeBenchmark TestType = "benchmark"
)

// TestConfig is a test configuration
type TestConfig struct {
	TestID     string
	Type       TestType
	Image      string
	Suite      string
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
