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

	"github.com/onosproject/onos-test/pkg/util/logging"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const defaultCluster = "onos"

// GetRootCommand returns the root onit command
func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "onit <command> [args]",
		Short:                  "Setup test clusters and run integration tests on Kubernetes",
		BashCompletionFunction: bashCompletion,
		SilenceUsage:           true,
	}
	cmd.AddCommand(getGetCommand())
	cmd.AddCommand(getCreateCommand())
	cmd.AddCommand(getDeleteCommand())
	cmd.AddCommand(getAddCommand())
	cmd.AddCommand(getRemoveCommand())
	cmd.AddCommand(getTestCommand())
	cmd.AddCommand(getBenchCommand())
	cmd.AddCommand(getSimulateCommand())
	cmd.AddCommand(getExecCommand())
	cmd.AddCommand(getCompletionCommand())
	cmd.PersistentFlags().StringSlice("config", []string{SimConfigPath}, "set required config files")
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

// getCluster returns the cluster from the given command flags
func getCluster(cmd *cobra.Command) string {
	cluster, _ := cmd.Flags().GetString("cluster")
	if cluster == "" {
		cluster = defaultCluster
	}
	return cluster
}

// getName returns the name of a resource from the given command flags
func getName(cmd *cobra.Command) string {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		name = random.NewPetName(2)
	}
	return name
}

func setupCommand(cmd *cobra.Command) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	logging.SetVerbose(verbose)
}
