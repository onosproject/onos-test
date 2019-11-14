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
	"k8s.io/api/core/v1"
)

// TestJob manages a single test job for a suite
type TestJob struct {
	cluster *TestCluster
	test    *TestConfig
}

// start starts the test job
func (j *TestJob) start() error {
	if err := j.cluster.Create(); err != nil {
		return err
	}
	if err := j.cluster.StartTest(j.test); err != nil {
		return err
	}
	if err := j.cluster.awaitTestJobRunning(j.test); err != nil {
		return err
	}
	return nil
}

// getStatus gets the status message and exit code of the given pod
func (j *TestJob) getStatus() (string, int, error) {
	return j.cluster.GetTestResult(j.test)
}

// getPod finds the Pod for the given test
func (j *TestJob) getPod() (*v1.Pod, error) {
	return j.cluster.getPod(j.test)
}

// tearDown tears down the job
func (j *TestJob) tearDown() error {
	return j.cluster.Delete()
}
