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

// Run runs the test job
func (j *TestJob) Run() (string, int, error) {
	if err := j.start(); err != nil {
		return "", 0, err
	}
	if err := j.waitForComplete(); err != nil {
		_ = j.tearDown()
		return "", 0, err
	}
	message, code, err := j.getResult()
	_ = j.tearDown()
	return message, code, err
}

// start starts the test job
func (j *TestJob) start() error {
	if err := j.cluster.Create(); err != nil {
		return err
	}
	if err := j.cluster.StartTest(j.test); err != nil {
		return err
	}
	return nil
}

// waitForComplete waits for the test job to finish running
func (j *TestJob) waitForComplete() error {
	if err := j.cluster.AwaitTestComplete(j.test); err != nil {
		return err
	}
	return nil
}

// getResult gets the job result
func (j *TestJob) getResult() (string, int, error) {
	return j.cluster.GetTestResult(j.test)
}

// tearDown tears down the job
func (j *TestJob) tearDown() error {
	return j.cluster.Delete()
}
