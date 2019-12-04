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

package grpc

import (
	"context"
	"github.com/onosproject/onos-test/pkg/benchmark"
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"google.golang.org/grpc"
)

// BenchmarkSuite :: benchmark
type BenchmarkSuite struct {
	benchmark.Suite
	client TestServiceClient
}

// SetupBenchmarkSuite :: benchmark
func (s *BenchmarkSuite) SetupBenchmarkSuite(c *benchmark.Context) {
	setup.App("test").
		SetImage("onosproject/grpc-test:latest").
		AddPort("grpc", 8080).
		SetReplicas(c.GetArg("replicas").Int(1))
	setup.SetupOrDie()
}

// SetupBenchmark :: benchmark
func (s *BenchmarkSuite) SetupBenchmark(b *benchmark.Benchmark) {
	conn, err := grpc.Dial("test:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	s.client = NewTestServiceClient(conn)
}

// BenchmarkGRPCRequestReply :: benchmark
func (s *BenchmarkSuite) BenchmarkGRPCRequestReply(b *benchmark.Benchmark) {
	params := []benchmark.Param{
		benchmark.RandomBytes(b.GetArg("value-count").Int(1), b.GetArg("value-length").Int(128)),
	}
	b.Run(func(value []byte) error {
		_, err := s.client.RequestReply(context.Background(), &Message{
			Value: value,
		})
		return err
	}, params...)
}
