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
	corev1 "k8s.io/api/core/v1"
	"os"
	"strconv"
	"strings"
)

type benchmarkContext string

const (
	benchmarkNamespaceEnv = "BENCHMARK_NAMESPACE"
	benchmarkContextEnv   = "BENCHMARK_CONTEXT"

	benchmarkJobEnv             = "BENCHMARK_JOB"
	benchmarkImageEnv           = "BENCHMARK_IMAGE"
	benchmarkImagePullPolicyEnv = "BENCHMARK_IMAGE_PULL_POLICY"
	benchmarkSuiteEnv           = "BENCHMARK_SUITE"
	benchmarkNameEnv            = "BENCHMARK_NAME"
	benchmarkWorkersEnv         = "BENCHMARK_WORKERS"
	benchmarkParallelismEnv     = "BENCHMARK_PARALLELISM"
	benchmarkRequestsEnv        = "BENCHMARK_REQUESTS"
	benchmarkArgPrefix          = "BENCHMARK_ARG_"
)

const (
	benchmarkContextCoordinator benchmarkContext = "coordinator"
	benchmarkContextWorker      benchmarkContext = "worker"
)

// getBenchmarkContext returns the current benchmark context
func getBenchmarkContext() benchmarkContext {
	context := os.Getenv(benchmarkContextEnv)
	if context != "" {
		return benchmarkContext(context)
	}
	return benchmarkContextCoordinator
}

func getBenchmarkEnv() map[string]string {
	env := make(map[string]string)
	for _, keyval := range os.Environ() {
		key := keyval[:strings.Index(keyval, "=")]
		value := keyval[strings.Index(keyval, "=")+1:]
		env[key] = value
	}
	return env
}

func getBenchmarkJob() string {
	return os.Getenv(benchmarkJobEnv)
}

func getBenchmarkNamespace() string {
	return os.Getenv(benchmarkNamespaceEnv)
}

func getBenchmarkImage() string {
	return os.Getenv(benchmarkImageEnv)
}

func getBenchmarkImagePullPolicy() corev1.PullPolicy {
	return corev1.PullPolicy(os.Getenv(benchmarkImagePullPolicyEnv))
}

func getBenchmarkSuite() string {
	return os.Getenv(benchmarkSuiteEnv)
}

func getBenchmarkName() string {
	return os.Getenv(benchmarkNameEnv)
}

func getBenchmarkWorkers() int {
	workers := os.Getenv(benchmarkWorkersEnv)
	if workers == "" {
		return 1
	}
	i, err := strconv.Atoi(workers)
	if err != nil {
		panic(err)
	}
	return i
}

func getBenchmarkParallelism() int {
	parallelism := os.Getenv(benchmarkParallelismEnv)
	if parallelism == "" {
		return 1
	}
	i, err := strconv.Atoi(parallelism)
	if err != nil {
		panic(err)
	}
	return i
}

func getBenchmarkRequests() int {
	requests := os.Getenv(benchmarkRequestsEnv)
	if requests == "" {
		return 1
	}
	i, err := strconv.Atoi(requests)
	if err != nil {
		panic(err)
	}
	return i
}

func getBenchmarkArgs() map[string]string {
	args := make(map[string]string)
	for key, value := range getBenchmarkEnv() {
		if strings.HasPrefix(key, benchmarkArgPrefix) {
			args[strings.ToLower(key[len(benchmarkArgPrefix):])] = value
		}
	}
	return args
}
