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
	"context"
	"fmt"
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/registry"
	"github.com/onosproject/onos-test/pkg/util/logging"
	"google.golang.org/grpc"
	"net"
	"reflect"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// newWorker returns a new benchmark worker
func newWorker(config *Config) (*Worker, error) {
	kubeAPI, err := kube.GetAPI(config.ID)
	if err != nil {
		return nil, err
	}
	return &Worker{
		client: kubeAPI.Client(),
		suites: make(map[string]BenchmarkingSuite),
	}, nil
}

// Worker runs a benchmark job
type Worker struct {
	client client.Client
	suites map[string]BenchmarkingSuite
}

// Run runs a benchmark
func (w *Worker) Run() error {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	RegisterWorkerServiceServer(server, w)
	return server.Serve(lis)
}

func (w *Worker) getSuite(name string) (BenchmarkingSuite, error) {
	if suite, ok := w.suites[name]; ok {
		return suite, nil
	}
	if suite := registry.GetBenchmarkSuite(name); suite != nil {
		w.suites[name] = suite
		return suite, nil
	}
	return nil, fmt.Errorf("unknown benchmark suite %s", name)
}

// SetupSuite sets up a benchmark suite
func (w *Worker) SetupSuite(ctx context.Context, request *SuiteRequest) (*SuiteResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Suite, getBenchmarkWorker()), "SetupSuite %s", request.Suite)
	step.Start()

	suite, err := w.getSuite(request.Suite)
	if err != nil {
		step.Fail(err)
		return nil, err
	}

	if setupSuite, ok := suite.(SetupSuite); ok {
		setupSuite.SetupSuite(newContext(request.Suite, request.Args))
	}

	step.Complete()
	return &SuiteResponse{}, nil
}

// TearDownSuite tears down a benchmark suite
func (w *Worker) TearDownSuite(ctx context.Context, request *SuiteRequest) (*SuiteResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Suite, getBenchmarkWorker()), "TearDownSuite %s", request.Suite)
	step.Start()

	suite, err := w.getSuite(request.Suite)
	if err != nil {
		step.Fail(err)
		return nil, err
	}

	if tearDownSuite, ok := suite.(TearDownSuite); ok {
		tearDownSuite.TearDownSuite(newContext(request.Suite, request.Args))
	}

	step.Complete()
	return &SuiteResponse{}, nil
}

// SetupWorker sets up a benchmark worker
func (w *Worker) SetupWorker(ctx context.Context, request *SuiteRequest) (*SuiteResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Suite, getBenchmarkWorker()), "SetupWorker %s", request.Suite)
	step.Start()

	suite, err := w.getSuite(request.Suite)
	if err != nil {
		step.Fail(err)
		return nil, err
	}

	if setupWorker, ok := suite.(SetupWorker); ok {
		setupWorker.SetupWorker(newContext(request.Suite, request.Args))
	}

	step.Complete()
	return &SuiteResponse{}, nil
}

// TearDownWorker tears down a benchmark worker
func (w *Worker) TearDownWorker(ctx context.Context, request *SuiteRequest) (*SuiteResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Suite, getBenchmarkWorker()), "TearDownWorker %s", request.Suite)
	step.Start()

	suite, err := w.getSuite(request.Suite)
	if err != nil {
		step.Fail(err)
		return nil, err
	}

	if tearDownWorker, ok := suite.(TearDownWorker); ok {
		tearDownWorker.TearDownWorker(newContext(request.Suite, request.Args))
	}

	step.Complete()
	return &SuiteResponse{}, nil
}

// SetupBenchmark sets up a benchmark
func (w *Worker) SetupBenchmark(ctx context.Context, request *BenchmarkRequest) (*BenchmarkResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Suite, getBenchmarkWorker()), "SetupBenchmark %s", request.Benchmark)
	step.Start()

	suite, err := w.getSuite(request.Suite)
	if err != nil {
		step.Fail(err)
		return nil, err
	}

	context := newContext(request.Benchmark, request.Args)
	if setupBenchmark, ok := suite.(SetupBenchmark); ok {
		setupBenchmark.SetupBenchmark(context)
	}

	methods := reflect.TypeOf(suite)
	if method, ok := methods.MethodByName("Setup" + request.Benchmark); ok {
		method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(context)})
	}

	step.Complete()
	return &BenchmarkResponse{}, nil
}

// TearDownBenchmark tears down a benchmark
func (w *Worker) TearDownBenchmark(ctx context.Context, request *BenchmarkRequest) (*BenchmarkResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Suite, getBenchmarkWorker()), "TearDownBenchmark %s", request.Benchmark)
	step.Start()

	suite, err := w.getSuite(request.Suite)
	if err != nil {
		step.Fail(err)
		return nil, err
	}

	context := newContext(request.Benchmark, request.Args)
	if tearDownBenchmark, ok := suite.(TearDownBenchmark); ok {
		tearDownBenchmark.TearDownBenchmark(context)
	}

	methods := reflect.TypeOf(suite)
	if method, ok := methods.MethodByName("TearDown" + request.Benchmark); ok {
		method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(context)})
	}

	step.Complete()
	return &BenchmarkResponse{}, nil
}

// RunBenchmark runs a benchmark
func (w *Worker) RunBenchmark(ctx context.Context, request *RunRequest) (*RunResponse, error) {
	step := logging.NewStep(fmt.Sprintf("%s/%d", request.Suite, getBenchmarkWorker()), "RunBenchmark %s", request.Benchmark)
	step.Start()

	suite, err := w.getSuite(request.Suite)
	if err != nil {
		step.Fail(err)
		return nil, err
	}

	methods := reflect.TypeOf(suite)
	method, ok := methods.MethodByName(request.Benchmark)
	if !ok {
		err = fmt.Errorf("unknown benchmark method %s", request.Benchmark)
		step.Fail(err)
		return nil, err
	}

	context := newContext(request.Benchmark, request.Args)
	b := newBenchmark(request.Benchmark, int(request.Requests), request.Duration, int(request.Parallelism), request.MaxLatency, context)
	method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(b)})

	step.Complete()
	return b.getResult(), nil
}

// benchmarkFilter filters benchmark method names
func benchmarkFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Benchmark", name); !ok {
		return false, nil
	}
	return true, nil
}
