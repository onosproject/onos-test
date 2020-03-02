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
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/util"
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
	benchmarkDurationEnv        = "BENCHMARK_DURATION"
	benchmarkArgsEnv            = "BENCHMARK_ARGS"
	benchmarkWorkerEnv          = "BENCHMARK_WORKER"
	benchmarkMaxLatencyMSEnv    = "BENCHMARK_MAX_LATENCY_MS"
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
	for key, value := range cluster.GetArgs() {
		args[key] = value
	}
	for key, value := range util.SplitMap(os.Getenv(benchmarkArgsEnv)) {
		args[key] = value
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
	var duration *time.Duration
	var maxLatency *time.Duration

	durationEnv := os.Getenv(benchmarkDurationEnv)
	if durationEnv != "" {
		d, err := strconv.Atoi(durationEnv)
		if err != nil {
			panic(err)
		}
		dur := time.Duration(d)
		duration = &dur
	}
	maxLatencyEnv := os.Getenv(benchmarkMaxLatencyMSEnv)
	if maxLatencyEnv != "" {
		d, err := strconv.Atoi(maxLatencyEnv)
		if err != nil {
			panic(err)
		}
		dur := time.Duration(d)
		maxLatency = &dur
	}

	return &Config{
		ID:              os.Getenv(benchmarkJobEnv),
		Image:           os.Getenv(benchmarkImageEnv),
		ImagePullPolicy: corev1.PullPolicy(os.Getenv(benchmarkImagePullPolicyEnv)),
		Suite:           os.Getenv(benchmarkSuiteEnv),
		Benchmark:       os.Getenv(benchmarkNameEnv),
		MaxLatency:      maxLatency,
		Workers:         workers,
		Parallelism:     parallelism,
		Requests:        requests,
		Duration:        duration,
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
	Duration        *time.Duration
	Args            map[string]string
	Env             map[string]string
	Timeout         time.Duration
	MaxLatency      *time.Duration
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
	if c.MaxLatency != nil {
		env[benchmarkMaxLatencyMSEnv] = fmt.Sprintf("%d", *c.MaxLatency)
	}
	if c.Duration != nil {
		env[benchmarkDurationEnv] = fmt.Sprintf("%d", *c.Duration)
	}
	env[benchmarkArgsEnv] = util.JoinMap(c.Args)
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
