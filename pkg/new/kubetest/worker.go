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
	"fmt"
	"github.com/onosproject/onos-test/pkg/new/util/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

// newTestWorker returns a new test worker
func newTestWorker(test *TestConfig) (Worker, error) {
	client, err := k8s.GetClient()
	if err != nil {
		return nil, err
	}
	return &TestWorker{
		client: client,
		test:   test,
	}, nil
}

// newBenchmarkWorker returns a new test worker
func newBenchmarkWorker(test *TestConfig) (Worker, error) {
	client, err := k8s.GetClient()
	if err != nil {
		return nil, err
	}
	return &BenchmarkWorker{
		client: client,
		test:   test,
	}, nil
}

// Worker runs a single test suite
type Worker interface {
	// Run runs a test suite
	Run() error
}

// TestWorker runs a test job
type TestWorker struct {
	client client.Client
	test   *TestConfig
}

// Run runs a test
func (w *TestWorker) Run() error {
	test := Registry.tests[w.test.Suite]
	if test == nil {
		return fmt.Errorf("unknown test suite %s", w.test.Suite)
	}

	tests := []testing.InternalTest{
		{
			Name: w.test.Suite,
			F: func(t *testing.T) {
				test.Run(t)
			},
		},
	}

	testing.Main(func(_, _ string) (bool, error) { return true, nil }, tests, nil, nil)
	return nil
}

// BenchmarkWorker runs a benchmark job
type BenchmarkWorker struct {
	client client.Client
	test   *TestConfig
}

// Run runs a benchmark
func (w *BenchmarkWorker) Run() error {
	benchmark := Registry.benchmarks[w.test.Suite]
	if benchmark == nil {
		return fmt.Errorf("unknown benchmark suite %s", w.test.Suite)
	}

	benchmarks := []testing.InternalBenchmark{
		{
			Name: w.test.Suite,
			F: func(b *testing.B) {
				benchmark.Run(b)
			},
		},
	}

	testing.Main(func(_, _ string) (bool, error) { return true, nil }, nil, benchmarks, nil)
	return nil
}
