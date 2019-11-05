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

// TestJob manages a single test job for a suite
type TestJob struct {
	cluster *TestCluster
	test    *TestConfig
}

// Start starts the test job
func (j *TestJob) Start() error {
	if err := j.cluster.Create(); err != nil {
		return err
	}
	if err := j.cluster.StartTest(j.test); err != nil {
		return err
	}
	return nil
}

// WaitForComplete waits for the test job to finish running
func (j *TestJob) WaitForComplete() error {
	if err := j.cluster.AwaitTestComplete(j.test); err != nil {
		return err
	}
	return nil
}

// GetResult gets the job result
func (j *TestJob) GetResult() (string, int, error) {
	return j.cluster.GetTestResult(j.test)
}

// TearDown tears down the job
func (j *TestJob) TearDown() error {
	return j.cluster.Delete()
}
