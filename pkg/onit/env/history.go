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

package env

import (
	"github.com/onosproject/onos-test/pkg/onit/cluster"
)

// HistoryEnv provides the history environment
type HistoryEnv interface {

	// GetTestsMap returns a map of tests
	GetTestsMap() map[string]cluster.JobInfo

	// GetBenchmarksMap returns a map of benchmarks
	GetBenchmarksMap() map[string]cluster.JobInfo

	// ListTests returns a history of tests
	ListTests() []cluster.JobInfo

	// ListBenchmarks returns a history of benchmarks
	ListBenchmarks() []cluster.JobInfo
}

var _ HistoryEnv = &clusterHistoryEnv{}

// clusterNetworksEnv is an implementation of the Networks interface
type clusterHistoryEnv struct {
	history *cluster.History
}

// ListTests gets list of tests
func (e *clusterHistoryEnv) ListTests() []cluster.JobInfo {
	return e.history.ListTests()
}

// ListBenchmarks gets list of benchmarks
func (e *clusterHistoryEnv) ListBenchmarks() []cluster.JobInfo {
	return e.history.ListBenchmarks()
}

// GetTestsMap gets a map of tests
func (e *clusterHistoryEnv) GetTestsMap() map[string]cluster.JobInfo {
	return e.history.GetTestsMap()
}

// GetBenchmarksMap gets a map of benchmarks
func (e *clusterHistoryEnv) GetBenchmarksMap() map[string]cluster.JobInfo {
	return e.history.GetBenchmarksMap()
}
