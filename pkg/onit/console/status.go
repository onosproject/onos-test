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

package console

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"time"
)

var (
	success = color.GreenString("✓")
	failure = color.RedString("✗")
)

const (
	waitDuration = 500 * time.Millisecond
)

type statusType string

const (
	statusStart    statusType = "start"
	statusProgress statusType = "progress"
	statusSucceed  statusType = "succeed"
	statusFail     statusType = "fail"
)

// ErrorStatus tracks the errors that occurred for a long-running operation
type ErrorStatus interface {
	// Failed returns a boolean indicating whether any failures have occurred
	Failed() bool

	// Errors returns a list of errors that occurred
	Errors() []error
}

// NewStatusWriter creates a new default StatusWriter
func NewStatusWriter() *StatusWriter {
	writer := os.Stdout
	spinner := newSpinner(writer)
	s := &StatusWriter{
		spinner:    spinner,
		writer:     writer,
		errors:     []error{},
		updateCh:   make(chan statusUpdate, 10),
		completeCh: make(chan bool, 1),
	}
	s.reset()
	return s
}

// StatusWriter provides real-time status output during onit setup operations
type StatusWriter struct {
	ErrorStatus
	spinner     *Spinner
	updateCh    chan statusUpdate
	status      string
	writer      io.Writer
	errors      []error
	lastUpdated *time.Time
	completeCh  chan bool
	complete    bool
}

// statusUpdate is an update the status
type statusUpdate struct {
	statusType statusType
	message    string
	err        error
}

// process starts processing status update events
func (s *StatusWriter) process() {
	for update := range s.updateCh {
		switch update.statusType {
		case statusStart:
			s.start(update.message)
		case statusProgress:
			s.progress(update.message)
		case statusSucceed:
			s.succeed()
		case statusFail:
			s.fail(update.err)
		}
	}
	s.completeCh <- true
}

// reset resets the state of the writer
func (s *StatusWriter) reset() {
	s.updateCh = make(chan statusUpdate)
	s.completeCh = make(chan bool)
	s.complete = false
	go s.process()
}

// wait waits a brief period if necessary to ensure status updates are not coming too quickly
func (s *StatusWriter) wait() {
	t := time.Now()
	if s.lastUpdated != nil && t.Sub(*s.lastUpdated) < waitDuration {
		time.Sleep((*s.lastUpdated).Add(waitDuration).Sub(t))
	}
	s.lastUpdated = &t
}

// Start starts a new status and begins a loading spinner
func (s *StatusWriter) Start(status string) *StatusWriter {
	if s.complete {
		s.reset()
	}
	s.updateCh <- statusUpdate{
		statusType: statusStart,
		message:    status,
	}
	return s
}

// start writes a new status and begins a loading spinner
func (s *StatusWriter) start(status string) {
	s.succeed()
	// set new status
	s.status = status
	s.spinner.SetMessage(fmt.Sprintf(" %s ", s.status))
	s.spinner.Spin()
}

// Progress appends a progress message to the current status
func (s *StatusWriter) Progress(message string) *StatusWriter {
	s.updateCh <- statusUpdate{
		statusType: statusProgress,
		message:    message,
	}
	return s
}

// Progress prints a progress message to the status
func (s *StatusWriter) progress(message string) {
	if s.status == "" {
		return
	}

	s.wait()
	s.spinner.SetMessage(fmt.Sprintf(" %-24s %-24s ", s.status, message))
}

// Succeed completes the current status successfully
func (s *StatusWriter) Succeed() *StatusWriter {
	s.updateCh <- statusUpdate{
		statusType: statusSucceed,
	}
	return s
}

// succeed removes the progress message and marks the current status as completed successfully
func (s *StatusWriter) succeed() {
	if s.status == "" {
		return
	}

	s.wait()
	s.spinner.Stop()
	fmt.Fprint(s.writer, "\r")
	fmt.Fprintf(s.writer, " %s %s\n", success, s.status)

	s.status = ""
}

// Fail fails the current status
func (s *StatusWriter) Fail(err error) *StatusWriter {
	s.errors = append(s.errors, err)
	s.updateCh <- statusUpdate{
		statusType: statusFail,
		err:        err,
	}
	return s
}

// fail marks the current status as failed and appends the error
func (s *StatusWriter) fail(err error) {
	if s.status == "" {
		return
	}

	s.wait()
	s.spinner.Stop()
	fmt.Fprint(s.writer, "\r")
	fmt.Fprintf(s.writer, " %s %-40s %s\n", failure, s.status, err)

	s.status = ""
}

// FlushAndClose flushes and closes the writer
func (s *StatusWriter) FlushAndClose() {
	if !s.complete {
		s.Close()
		<-s.completeCh
		s.complete = true
	}
}

// Failed returns a boolean indicating whether errors occurred
func (s *StatusWriter) Failed() bool {
	s.FlushAndClose()
	return len(s.errors) > 0
}

// Errors returns a list of errors that occurred
func (s *StatusWriter) Errors() []error {
	s.FlushAndClose()
	return s.errors
}

// Close closes the status writer
func (s *StatusWriter) Close() {
	close(s.updateCh)
}
