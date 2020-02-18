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

package model

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
)

var traceFile = path.Join(DataPath, "trace.json")

// NewService returns a new model checker service
func NewService() *Service {
	return &Service{}
}

// Service is a model checker service
type Service struct {
}

// Start starts the server
func (s *Service) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", CheckerPort))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	RegisterModelCheckerServiceServer(server, &modelCheckerServer{})
	if err := server.Serve(lis); err != nil {
		return err
	}
	return nil
}

// modelCheckerServer is a model checker service server
type modelCheckerServer struct{}

func (s *modelCheckerServer) CheckModel(request *ModelCheckRequest, stream ModelCheckerService_CheckModelServer) error {
	model := NewModel(request.Model)
	writer := &streamWriter{
		stream: stream,
	}
	if err := s.prepareTraces(model, writer); err != nil {
		return err
	}
	if err := s.runChecker(model, writer); err != nil {
		return err
	}
	return nil
}

func (s *modelCheckerServer) prepareTraces(model *Model, stream io.Writer) error {
	return os.Rename(model.dataPath, traceFile)
}

func (s *modelCheckerServer) runChecker(model *Model, stream io.Writer) error {
	cmd := exec.Command("tlc", "-metadir", path.Join(DataPath, "states"), model.specPath)
	cmd.Stdout = stream
	cmd.Stderr = stream
	return cmd.Run()
}

type streamWriter struct {
	stream ModelCheckerService_CheckModelServer
}

func (w *streamWriter) Write(bytes []byte) (int, error) {
	err := w.stream.Send(&ModelCheckResponse{
		State:   ModelCheckerState_RUNNING,
		Message: string(bytes),
	})
	if err != nil {
		return 0, err
	}
	return len(bytes), nil
}
