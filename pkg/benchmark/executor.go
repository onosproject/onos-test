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
	"os"
)

// The executor is the entrypoint for benchmark images. It takes the input and environment and runs
// the image in the appropriate context according to the arguments.

// benchmarkContext is the context in which a test image is running
type benchmarkContext string

const (
	testNamespaceEnv = "TEST_NAMESPACE"
	testContextEnv   = "TEST_CONTEXT"

	testContextCoordinator benchmarkContext = "coordinator"
	testContextWorker      benchmarkContext = "worker"
)

// Main runs a test
func Main() {
	if err := Run(); err != nil {
		println("Benchmark failed " + err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

// Run runs a test
func Run() error {
	context := getTestContext()
	switch context {
	case testContextCoordinator:
		config, err := loadCoordinatorConfig()
		if err != nil {
			return err
		}
		return runCoordinator(config)
	case testContextWorker:
		config, err := loadWorkerConfig()
		if err != nil {
			return err
		}
		return runWorker(config)
	}
	return nil
}

// getTestContext returns the current test context
func getTestContext() benchmarkContext {
	return benchmarkContext(os.Getenv(testContextEnv))
}

// runCoordinator runs a test image in the coordinator context
func runCoordinator(config *CoordinatorConfig) error {
	coordinator, err := newCoordinator(config)
	if err != nil {
		return err
	}
	return coordinator.Run()
}

// runWorker runs a test image in the worker context
func runWorker(config *WorkerConfig) error {
	worker, err := newWorker(config)
	if err != nil {
		return err
	}
	return worker.Run()
}
