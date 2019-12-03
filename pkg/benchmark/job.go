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

// Job manages a single benchmark job for a suite
type Job struct {
	cluster *Cluster
	config  *CoordinatorConfig
}

// start starts the test job
func (j *Job) run() error {
	if err := j.cluster.Create(); err != nil {
		return err
	}
	if err := j.cluster.CreateWorkers(j.config); err != nil {
		return err
	}
	if err := j.cluster.RunBenchmarks(j.config); err != nil {
		return err
	}
	return nil
}

// getStatus gets the status message and exit code of the given pod
func (j *Job) getStatus() (string, int, error) {
	return j.cluster.GetResult(j.config)
}

// tearDown tears down the job
func (j *Job) tearDown() error {
	return j.cluster.Delete()
}
