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

package benchmark

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"os"
	"strconv"
	"strings"
	"time"
)

type benchmarkContext string

const (
	benchmarkContextEnv = "BENCHMARK_CONTEXT"

	benchmarkJobEnv             = "BENCHMARK_JOB"
	benchmarkImageEnv           = "BENCHMARK_IMAGE"
	benchmarkImagePullPolicyEnv = "BENCHMARK_IMAGE_PULL_POLICY"
	benchmarkSuiteEnv           = "BENCHMARK_SUITE"
	benchmarkNameEnv            = "BENCHMARK_NAME"
	benchmarkWorkersEnv         = "BENCHMARK_WORKERS"
	benchmarkParallelismEnv     = "BENCHMARK_PARALLELISM"
	benchmarkRequestsEnv        = "BENCHMARK_REQUESTS"
	benchmarkArgPrefix          = "BENCHMARK_ARG_"
	benchmarkWorkerEnv          = "BENCHMARK_WORKER"
)

const (
	benchmarkContextCoordinator benchmarkContext = "coordinator"
	benchmarkContextWorker      benchmarkContext = "worker"
)

// GetConfigFromEnv returns the benchmark configuration from the environment
func GetConfigFromEnv() *Config {
	env := make(map[string]string)
	for _, keyval := range os.Environ() {
		key := keyval[:strings.Index(keyval, "=")]
		value := keyval[strings.Index(keyval, "=")+1:]
		env[key] = value
	}
	args := make(map[string]string)
	for key, value := range env {
		if strings.HasPrefix(key, benchmarkArgPrefix) {
			args[strings.ToLower(key[len(benchmarkArgPrefix):])] = value
		}
	}
	workers, err := strconv.Atoi(os.Getenv(benchmarkWorkersEnv))
	if err != nil {
		panic(err)
	}
	parallelism, err := strconv.Atoi(os.Getenv(benchmarkParallelismEnv))
	if err != nil {
		panic(err)
	}
	requests, err := strconv.Atoi(os.Getenv(benchmarkRequestsEnv))
	if err != nil {
		panic(err)
	}
	return &Config{
		ID:              os.Getenv(benchmarkJobEnv),
		Image:           os.Getenv(benchmarkImageEnv),
		ImagePullPolicy: corev1.PullPolicy(os.Getenv(benchmarkImagePullPolicyEnv)),
		Suite:           os.Getenv(benchmarkSuiteEnv),
		Benchmark:       os.Getenv(benchmarkNameEnv),
		Workers:         workers,
		Parallelism:     parallelism,
		Requests:        requests,
		Args:            args,
		Env:             env,
	}
}

// Config is a benchmark configuration
type Config struct {
	ID              string
	Image           string
	ImagePullPolicy corev1.PullPolicy
	Suite           string
	Benchmark       string
	Workers         int
	Parallelism     int
	Requests        int
	Args            map[string]string
	Env             map[string]string
	Timeout         time.Duration
}

// ToEnv returns the configuration as a mapping of environment variables
func (c *Config) ToEnv() map[string]string {
	env := c.Env
	env[benchmarkJobEnv] = c.ID
	env[benchmarkImageEnv] = c.Image
	env[benchmarkImagePullPolicyEnv] = string(c.ImagePullPolicy)
	env[benchmarkSuiteEnv] = c.Suite
	env[benchmarkNameEnv] = c.Benchmark
	env[benchmarkWorkersEnv] = fmt.Sprintf("%d", c.Workers)
	env[benchmarkParallelismEnv] = fmt.Sprintf("%d", c.Parallelism)
	env[benchmarkRequestsEnv] = fmt.Sprintf("%d", c.Requests)
	for key, value := range c.Args {
		env[benchmarkArgPrefix+strings.ToUpper(key)] = value
	}
	return env
}

// getBenchmarkContext returns the current benchmark context
func getBenchmarkContext() benchmarkContext {
	context := os.Getenv(benchmarkContextEnv)
	if context != "" {
		return benchmarkContext(context)
	}
	return benchmarkContextCoordinator
}

// getBenchmarkWorker returns the current benchmark worker number
func getBenchmarkWorker() int {
	worker := os.Getenv(benchmarkWorkerEnv)
	if worker == "" {
		return 0
	}
	i, err := strconv.Atoi(worker)
	if err != nil {
		panic(err)
	}
	return i
}
