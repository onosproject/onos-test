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
	"github.com/onosproject/onos-test/pkg/onit"
	"github.com/onosproject/onos-test/pkg/onit/benchmark"
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"google.golang.org/grpc"
	"k8s.io/api/core/v1"
)

type GRPCBenchmarkSuite struct {
	onit.ScriptSuite
}

func (b *GRPCBenchmarkSuite) SetupScriptSuite() {
	setup.App("test").
		SetImage("192.168.1.11:30000/onosproject/grpc-test:latest").
		SetPullPolicy(v1.PullAlways).
		AddPort("grpc", 8080).
		SetReplicas(1)
	setup.SetupOrDie()
}

func (b *GRPCBenchmarkSuite) getClient() TestServiceClient {
	conn, err := grpc.Dial("test:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return NewTestServiceClient(conn)
}

func (b *GRPCBenchmarkSuite) RunBenchmarkRequests() {
	benchmark.New().
		SetHandlerFactory(func() benchmark.Handler {
			return &GRPCBenchmarkRequestHandler{
				client: b.getClient(),
			}
		}).
		SetClients(1).
		SetParallelism(1).
		SetRequests(1).
		AddHandlerArg(benchmark.RandomBytes(1000, 128)).
		Run()
}

type GRPCBenchmarkRequestHandler struct {
	client TestServiceClient
}

func (h *GRPCBenchmarkRequestHandler) Run(args ...interface{}) error {
	_, err := h.client.RequestReply(context.Background(), &Message{
		Value: args[0].([]byte),
	})
	return err
}
