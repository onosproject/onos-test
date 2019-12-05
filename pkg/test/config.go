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
	"k8s.io/api/core/v1"
	"os"
	"strings"
)

type testContext string

const (
	testNamespaceEnv = "TEST_NAMESPACE"
	testContextEnv   = "TEST_CONTEXT"

	testJobEnv             = "TEST_JOB"
	testImageEnv           = "TEST_IMAGE"
	testImagePullPolicyEnv = "TEST_IMAGE_PULL_POLICY"
	testSuiteEnv           = "TEST_SUITE"
	testNameEnv            = "TEST_NAME"
)

const (
	testContextCoordinator testContext = "coordinator"
	testContextWorker      testContext = "worker"
)

// getTestContext returns the current test context
func getTestContext() testContext {
	context := os.Getenv(testContextEnv)
	if context != "" {
		return testContext(context)
	}
	return testContextCoordinator
}

func getTestEnv() map[string]string {
	env := make(map[string]string)
	for _, keyval := range os.Environ() {
		key := keyval[:strings.Index(keyval, "=")]
		value := keyval[strings.Index(keyval, "=")+1:]
		env[key] = value
	}
	return env
}

func getTestJob() string {
	return os.Getenv(testJobEnv)
}

func getTestNamespace() string {
	return os.Getenv(testNamespaceEnv)
}

func getTestImage() string {
	return os.Getenv(testImageEnv)
}

func getTestImagePullPolicy() v1.PullPolicy {
	return v1.PullPolicy(os.Getenv(testImagePullPolicyEnv))
}

func getTestSuite() string {
	return os.Getenv(testSuiteEnv)
}

func getTestName() string {
	return os.Getenv(testNameEnv)
}
