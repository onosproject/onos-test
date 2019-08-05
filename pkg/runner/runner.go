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

package runner

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

// TestRunner runs integration tests
type TestRunner struct {
	Registry *TestRegistry
}

// RunTests Runs the tests
func (r *TestRunner) RunTests(args []string) error {
	tests := make([]testing.InternalTest, 0, len(args))
	if len(args) > 0 {
		for _, name := range args {
			test, ok := r.Registry.tests[name]
			if !ok {
				return errors.New("unknown test " + name)
			}
			tests = append(tests, testing.InternalTest{
				Name: name,
				F:    test,
			})
		}
	} else {
		for name, test := range r.Registry.tests {
			tests = append(tests, testing.InternalTest{
				Name: name,
				F:    test,
			})
		}
	}

	// Hack to enable verbose testing.
	os.Args = []string{
		os.Args[0],
		"-test.v",
	}

	// Run the integration tests via the testing package.
	testing.Main(func(_, _ string) (bool, error) { return true, nil }, tests, nil, nil)
	return nil
}

// RunTestSuites Runs the tests groups
func (r *TestRunner) RunTestSuites(args []string) error {
	for _, name := range args {
		testSuite, ok := r.Registry.TestSuites[name]
		if !ok {
			return errors.New("unknown test suite" + name)
		}
		err := r.RunTests(testSuite.GetTestNames())
		if err != nil {
			return err
		}
	}
	return nil
}

// RunBenchmarks runs the benchmarks
func (r *TestRunner) RunBenchmarks(args []string, n int) error {
	benchmarks := make([]testing.InternalBenchmark, 0, len(args))
	if len(args) > 0 {
		for _, name := range args {
			benchmark, ok := r.Registry.benchmarks[name]
			if !ok {
				return errors.New("unknown benchmark " + name)
			}
			benchmarks = append(benchmarks, testing.InternalBenchmark{
				Name: name,
				F:    benchmark,
			})
		}
	} else {
		for name, benchmark := range r.Registry.benchmarks {
			benchmarks = append(benchmarks, testing.InternalBenchmark{
				Name: name,
				F:    benchmark,
			})
		}
	}

	// Hack to enable verbose testing.
	os.Args = []string{
		os.Args[0],
		"-test.bench=.",
		"-test.v",
	}

	// If a count was specified, append the count.
	if n > 0 {
		os.Args = append(os.Args, fmt.Sprintf("-test.count=%d", n))
	}

	// Run the integration tests via the testing package.
	testing.Main(func(_, _ string) (bool, error) { return true, nil }, nil, benchmarks, nil)
	return nil
}

// RunBenchmarkSuites Runs a benchmark suite
func (r *TestRunner) RunBenchmarkSuites(args []string, n int) error {
	for _, name := range args {
		benchSuite, ok := r.Registry.BenchSuites[name]
		if !ok {
			return errors.New("unknown test suite" + name)
		}
		err := r.RunBenchmarks(benchSuite.GetBenchNames(), n)
		if err != nil {
			return err
		}
	}
	return nil
}
