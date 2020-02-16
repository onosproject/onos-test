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
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"os"
)

// Checker is a model checker
type Checker interface {
	// CheckModel checks the given model against the given traces
	CheckModel(model *Model, traces []*Trace) error
}

// NewChecker gets a model checker client
func NewChecker() (Checker, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	return &checkerClient{
		client: client,
	}, nil
}

// newClient creates a model checker client
func newClient() (ModelCheckerServiceClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", CheckerPort), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return NewModelCheckerServiceClient(conn), nil
}

// checkerClient is a model checker client
type checkerClient struct {
	client ModelCheckerServiceClient
}

func (c *checkerClient) CheckModel(model *Model, traces []*Trace) error {
	if err := c.writeTraces(model, traces); err != nil {
		return err
	}
	if err := c.callCheck(model); err != nil {
		return err
	}
	return nil
}

// writeTrace writes the trace
func (c *checkerClient) writeTraces(model *Model, traces []*Trace) error {
	file, err := os.Create(model.dataPath)
	if err != nil {
		return err
	}
	for _, trace := range traces {
		if _, err := fmt.Fprintln(file, string(trace.Bytes)); err != nil {
			return err
		}
	}
	return nil
}

// callCheck calls the model checker
func (c *checkerClient) callCheck(model *Model) error {
	request := &ModelCheckRequest{
		Model: model.Name,
	}

	stream, err := c.client.CheckModel(context.Background(), request)
	if err != nil {
		return err
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		switch response.State {
		case ModelCheckerState_RUNNING:
			fmt.Println(response.Message)
		case ModelCheckerState_PASSED:
			return nil
		case ModelCheckerState_FAILED:
			return errors.New("model validation failed")
		}
	}
}
