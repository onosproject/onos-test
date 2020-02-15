// Copyright 2020-present Open Networking Foundation.
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

package simulation

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-test/pkg/model"
	"google.golang.org/grpc"
	"net"
)

// Register records simulation traces
type Register interface {
	// Trace records a trace to the register
	Trace(values ...interface{}) error

	// close closes the register
	close()
}

// newBlockingRegister returns a new register that synchronously writes traces
func newBlockingRegister(address string) (Register, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &registerClient{
		conn: conn,
	}, nil
}

// registerClient is a register that synchronously writes traces to the register service
type registerClient struct {
	conn *grpc.ClientConn
}

func (r *registerClient) Trace(values ...interface{}) error {
	trace, err := model.NewTrace(values...)
	if err != nil {
		return err
	}

	client := NewRegisterServiceClient(r.conn)
	request := &TraceRequest{
		Trace: trace,
	}
	_, err = client.Trace(context.Background(), request)
	return err
}

func (r *registerClient) close() {
	_ = r.conn.Close()
}

// newRegisterServer returns a new register server
func newRegisterServer(address string) *registerServer {
	return &registerServer{
		address:   address,
		registers: make(map[string]chan<- *model.Trace),
	}
}

// registerServer is a server that listens for register writes
type registerServer struct {
	address   string
	server    *grpc.Server
	registers map[string]chan<- *model.Trace
}

// addRegister adds a register to the server
func (s *registerServer) addRegister(simulation string, ch chan<- *model.Trace) {
	s.registers[simulation] = ch
}

// serve starts serving the register
func (s *registerServer) serve() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	s.server = grpc.NewServer()
	RegisterRegisterServiceServer(s.server, s)
	return s.server.Serve(lis)
}

// Trace handles a register trace request
func (s *registerServer) Trace(ctx context.Context, request *TraceRequest) (*TraceResponse, error) {
	if request.Trace != nil {
		register, ok := s.registers[request.Simulation]
		if !ok {
			return nil, fmt.Errorf("unknown simulation %s", request.Simulation)
		}
		register <- request.Trace
	}
	return &TraceResponse{}, nil
}

// stop stops the server
func (s *registerServer) stop() {
	s.server.Stop()
}
