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
	"time"

	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/onosproject/onos-test/pkg/onit"
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/spf13/cobra"
)

var (
	runExample = templates.Examples(i18n.T(`
    # To run all integration tests:
    onit run suite

    # To run a single test on a cluster
    onit run test <name of a test>

`))
)

// getRunCommand returns a cobra run command to run integration tests
func getRunCommand(registry *runner.TestRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run {test,suite}",
		Short:   "Run integration tests",
		Example: runExample,
	}
	cmd.AddCommand(getRunSuiteRemoteCommand(registry))
	cmd.AddCommand(getRunTestRemoteCommand())
	return cmd
}

// getRunCommand returns a cobra "run" command
func getRunTestRemoteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test [tests]",
		Short: "Run integration tests on Kubernetes",
		Run: func(cmd *cobra.Command, args []string) {
			runTestsRemote(cmd, "test", args)
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster on which to run the test")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().IntP("timeout", "t", 60*10, "test timeout in seconds")
	return cmd
}

func getRunSuiteRemoteCommand(registry *runner.TestRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suite [suite]",
		Short: "Run integration tests",
		Run: func(cmd *cobra.Command, args []string) {
			runTestsRemote(cmd, "suite", args)
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster on which to run the test")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().IntP("timeout", "t", 60*10, "test timeout in seconds")

	return cmd
}

func runTestsRemote(cmd *cobra.Command, commandType string, tests []string) {
	testID := fmt.Sprintf("test-%d", newUUIDInt())

	// Get the onit controller
	controller, err := onit.NewController()
	if err != nil {
		exitError(err)
	}

	// Get the cluster ID
	clusterID, err := cmd.Flags().GetString("cluster")
	if err != nil {
		exitError(err)
	}

	// Get the cluster controller
	cluster, err := controller.GetCluster(clusterID)
	if err != nil {
		exitError(err)
	}

	timeout, _ := cmd.Flags().GetInt("timeout")
	message, code, status := cluster.RunTests(testID, append([]string{commandType}, tests...), time.Duration(timeout)*time.Second)
	if status.Failed() {
		exitStatus(status)
	} else {
		fmt.Println(message)
		os.Exit(code)
	}
}
