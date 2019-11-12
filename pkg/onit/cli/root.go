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
	"github.com/onosproject/onos-test/pkg/util/logging"
	"github.com/spf13/cobra"
)

// GetRootCommand returns the root onit command
func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "onit <command> [args]",
		Short:                  "Setup test clusters and run integration tests on Kubernetes",
		BashCompletionFunction: bashCompletion,
	}
	cmd.AddCommand(getGetCommand())
	cmd.AddCommand(getCreateCommand())
	cmd.AddCommand(getDeleteCommand())
	cmd.AddCommand(getAddCommand())
	cmd.AddCommand(getRemoveCommand())
	cmd.AddCommand(getRunCommand())
	cmd.AddCommand(getExecCommand())
	cmd.AddCommand(getCompletionCommand())
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	return cmd
}

func runCommand(cmd *cobra.Command) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	logging.SetVerbose(verbose)
}
