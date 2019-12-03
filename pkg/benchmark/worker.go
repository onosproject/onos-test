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
	"errors"
	"github.com/onosproject/onos-test/pkg/kube"
	"google.golang.org/grpc"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// newWorker returns a new benchmark worker
func newWorker(config *WorkerConfig) (*Worker, error) {
	kubeAPI, err := kube.GetAPI(config.JobID)
	if err != nil {
		return nil, err
	}
	suite := Registry.benchmarks[config.Suite]
	if suite == nil {
		return nil, errors.New("unknown benchmark suite")
	}
	return &Worker{
		client: kubeAPI.Client(),
		config: config,
		suite:  suite,
	}, nil
}

// Worker runs a benchmark job
type Worker struct {
	client client.Client
	config *WorkerConfig
	suite  BenchmarkingSuite
}

// Run runs a benchmark
func (w *Worker) Run() error {
	setupSuite(w.suite, w.config)
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	RegisterWorkerServiceServer(server, w)
	return server.Serve(lis)
}

func (w *Worker) RunBenchmark(ctx context.Context, request *Request) (*Result, error) {
	return runBenchmark(request.Benchmark, int(request.Requests), w.suite, w.config)
}
