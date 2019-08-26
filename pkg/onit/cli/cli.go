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

package cli

import (
	"fmt"
	"os"

	"github.com/onosproject/onos-test/pkg/runner"

	"github.com/google/uuid"
	"github.com/onosproject/onos-test/pkg/onit/console"
	"github.com/spf13/cobra"
)

// Contains tells whether array contains x.
func Contains(array []string, elem string) bool {
	for _, n := range array {
		if elem == n {
			return true
		}
	}
	return false
}

// Subset returns true if the first array is completely
// contained in the second array.
func Subset(first, second []string) bool {
	set := make(map[string]int)
	for _, value := range second {
		set[value]++
	}

	for _, value := range first {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}

	return true
}

// GetOnitCommand returns a Cobra command for tests in the given test registry
func GetOnitCommand(registry *runner.TestRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "onit",
		Short:                  "Run onos integration tests on Kubernetes",
		BashCompletionFunction: bashCompletion,
	}
	cmd.AddCommand(getCreateCommand())
	cmd.AddCommand(getAddCommand())
	cmd.AddCommand(getRemoveCommand())
	cmd.AddCommand(getDeleteCommand())
	cmd.AddCommand(getRunCommand(registry))
	cmd.AddCommand(getGetCommand(registry))
	cmd.AddCommand(getSetCommand())
	cmd.AddCommand(getDebugCommand())
	cmd.AddCommand(getFetchCommand())
	cmd.AddCommand(getCompletionCommand())
	cmd.AddCommand(getSSHCommand())
	cmd.AddCommand(getOnosCliCommand())

	return cmd
}

// newUUIDString returns a new string UUID
func newUUIDString() string {
	id, err := uuid.NewUUID()
	if err != nil {
		exitError(err)
	}
	return id.String()
}

// newUuidInt returns a numeric UUID
func newUUIDInt() uint32 {
	id, err := uuid.NewUUID()
	if err != nil {
		exitError(err)
	}
	return id.ID()
}

// exitStatus prints the errors from the given status and exits
func exitStatus(status console.ErrorStatus) {
	for _, err := range status.Errors() {
		fmt.Println(err)
	}
	os.Exit(1)
}

// exitError prints the given errors to stdout and exits with exit code 1
func exitError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
