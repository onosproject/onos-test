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

package cluster

import (
	v1 "k8s.io/api/core/v1"
)

func newHistory(cluster *Cluster) *History {
	return &History{
		client:  cluster.client,
		cluster: cluster,
	}
}

// New returns a new history
func (h *History) New() *History {
	return newHistory(h.cluster)
}

// History provides methods for retrieving history of tests and benchmarks
type History struct {
	*client
	cluster *Cluster
}

// JobInfo k8s job info
type JobInfo struct {
	jobName string
	status  string
	jobType string
	image   string
	envVar  map[string]string
}

// Status job status
type Status int

const (
	// Succeeded status
	Succeeded Status = iota
	// Failed status
	Failed
	// Running status
	Running
	// Pending status
	Pending
	// Unknown status
	Unknown
)

func (s Status) String() string {
	return [...]string{"Succeeded", "Failed", "Running", "Pending", "Unknown"}[s]
}

// GetEnvVar gets the job environment variables
func (j *JobInfo) GetEnvVar() map[string]string {
	return j.envVar
}

// GetJobImage gets job image
func (j *JobInfo) GetJobImage() string {
	return j.image
}

// GetJobName gets job name
func (j *JobInfo) GetJobName() string {
	return j.jobName
}

// GetJobStatus gets job status
func (j *JobInfo) GetJobStatus() string {
	return j.status
}

// GetJobType gets job type
func (j *JobInfo) GetJobType() string {
	return j.jobType
}

// getJobs gets list of all jobs based on a given jobType
func (h *History) getJobs(jobType string) ([]JobInfo, map[string]JobInfo) {
	pods := h.getPods(getLabels(jobType), "kube-test")
	var jobs []JobInfo
	jobsMap := make(map[string]JobInfo, len(pods.Items))

	for _, pod := range pods.Items {
		status := ""
		switch pod.Status.Phase {
		case v1.PodSucceeded:
			status = Succeeded.String()
		case v1.PodRunning:
			status = Running.String()

		case v1.PodFailed:
			status = Failed.String()

		case v1.PodPending:
			status = Pending.String()

		case v1.PodUnknown:
			status = Unknown.String()

		}
		envVars := pod.Spec.Containers[0].Env
		envVarMap := make(map[string]string, len(envVars))
		for _, envVar := range envVars {
			envVarMap[envVar.Name] = envVar.Value
		}

		job := JobInfo{
			jobName: pod.Name,
			status:  status,
			jobType: pod.Labels["type"],
			image:   pod.Spec.Containers[0].Image,
			envVar:  envVarMap,
		}
		jobs = append(jobs, job)
		jobsMap[job.jobName] = job

	}

	return jobs, jobsMap

}

// GetTestsMap gets a map of test jobs
func (h *History) GetTestsMap() map[string]JobInfo {
	_, testsMap := h.getJobs("test")
	return testsMap
}

// GetBenchmarksMap gets a maps of benchmarks
func (h *History) GetBenchmarksMap() map[string]JobInfo {
	_, benchmarksMap := h.getJobs("benchmark")
	return benchmarksMap
}

// ListTests gets list of all tests
func (h *History) ListTests() []JobInfo {
	testJobs, _ := h.getJobs("test")
	return testJobs
}

// ListBenchmarks gets list of all benchmarks
func (h *History) ListBenchmarks() []JobInfo {
	benchJobs, _ := h.getJobs("benchmark")
	return benchJobs
}
