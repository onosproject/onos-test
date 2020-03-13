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
	"math/rand"
	"os"
	"time"

	"github.com/onosproject/onos-test/pkg/util/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// GetRootCommand returns the root onit command
func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "onit <command> [args]",
		Short:                  "Setup test clusters and run integration tests on Kubernetes",
		BashCompletionFunction: bashCompletion,
		SilenceUsage:           true,
	}
	cmd.AddCommand(getTestCommand())
	cmd.AddCommand(getBenchCommand())
	cmd.AddCommand(getSimulateCommand())
	cmd.AddCommand(getCompletionCommand())
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	cmd.PersistentFlags().StringP("cluster", "c", "", "the cluster on which to execute the command")
	cmd.PersistentFlags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}

// GenerateCliDocs generate markdown files for onit commands
func GenerateCliDocs() {
	cmd := GetRootCommand()
	err := doc.GenMarkdownTree(cmd, "docs/cli")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func setupCommand(cmd *cobra.Command) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	logging.SetVerbose(verbose)
}
