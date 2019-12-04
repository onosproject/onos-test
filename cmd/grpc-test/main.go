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

package main

import (
	"context"
	grpc_bench "github.com/onosproject/onos-test/benchmark/grpc"
	"google.golang.org/grpc"
	"net"
)

func main() {
	service := &testService{}
	server := grpc.NewServer()
	grpc_bench.RegisterTestServiceServer(server, service)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}

type testService struct {
}

func (s *testService) RequestReply(ctx context.Context, message *grpc_bench.Message) (*grpc_bench.Message, error) {
	return message, nil
}
