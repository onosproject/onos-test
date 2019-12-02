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
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"github.com/onosproject/onos-test/pkg/test"
	"google.golang.org/grpc"
)

type GRPCBenchmarkSuite struct {
	test.BenchmarkSuite
}

func (s *GRPCBenchmarkSuite) SetupBenchmarkSuite(c *test.Context) {
	setup.App("test").
		SetImage("onosproject/grpc-test:latest").
		AddPort("grpc", 8080).
		SetReplicas(1)
	setup.SetupOrDie()
}

func (s *GRPCBenchmarkSuite) SetupBenchmark(b *test.Benchmark) {
	b.Init(func() (TestServiceClient, error) {
		conn, err := grpc.Dial("test:8080", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		return NewTestServiceClient(conn), nil
	})
}

func (s *GRPCBenchmarkSuite) BenchmarkGRPCRequestReply(b *test.Benchmark) {
	params := []test.Param{
		test.RandomBytes(b.GetArg("value-count").Int(), b.GetArg("value-length").Int()),
	}
	b.Run(func(client TestServiceClient, value []byte) error {
		_, err := client.RequestReply(context.Background(), &Message{
			Value: value,
		})
		return err
	}, params...)
}
