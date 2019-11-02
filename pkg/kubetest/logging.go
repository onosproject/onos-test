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

package kubetest

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"time"
)

var (
	success = "✓"
	failure = "✗"
	writer  = os.Stdout
)

// NewStep returns a new step
func NewStep(test, name string) *Step {
	return &Step{
		test: test,
		name: name,
	}
}

// Step is a loggable step
type Step struct {
	test string
	name string
}

// Start starts the step
func (s *Step) Start() {
	fmt.Fprintln(writer, fmt.Sprintf("  %s %s %s", time.Now().Format(time.RFC3339), s.test, s.name))
}

// Complete completes the step
func (s *Step) Complete() {
	fmt.Fprintln(writer, color.GreenString(fmt.Sprintf("%s %s %s %s", success, time.Now().Format(time.RFC3339), s.test, s.name)))
}

// Fail fails the step with the given error
func (s *Step) Fail(err error) {
	fmt.Fprintln(writer, color.RedString(fmt.Sprintf("%s %s %s %s", failure, time.Now().Format(time.RFC3339), s.test, s.name)))
}

// printLog prints the given log line
func printLog(line string) {
	if line[0] == success[0] {
		fmt.Fprintln(writer, color.GreenString(line))
	} else if line[0] == failure[0] {
		fmt.Fprintln(writer, color.RedString(line))
	} else {
		fmt.Fprintln(writer, line)
	}
}
