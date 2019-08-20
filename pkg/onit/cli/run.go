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

	"github.com/onosproject/onos-test/pkg/onit/k8s"

	"github.com/onosproject/onos-test/test"

	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/spf13/cobra"
)

var (
	runExample = `
		# Run a single test on the cluster
		onit run test <name of a test>

		# Run a suite of tests on the cluster
		onit run test-suite <name of a suite>

		# Run a benchmark on the cluster
		onit run bench <name of a benchmark>

		# Run a suite of benchmarks on the cluster
		onit run bench-suite <name of a suite>`
)

// getRunCommand returns a cobra run command to run integration tests
func getRunCommand(registry *runner.TestRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run {test,test-suite,bench,bench-suite}",
		Short:   "Run integration tests",
		Example: runExample,
	}
	cmd.AddCommand(getRunTestCommand())
	cmd.AddCommand(getRunTestSuiteCommand(registry))
	cmd.AddCommand(getRunBenchCommand())
	cmd.AddCommand(getRunBenchSuiteCommand(registry))
	return cmd
}

func getRunTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test [tests]",
		Short: "Run tests on Kubernetes",
		Run: func(cmd *cobra.Command, args []string) {
			count, _ := cmd.Flags().GetInt("count")
			testNames := test.Registry.GetTestNames()
			testID := fmt.Sprintf("test-%d", newUUIDInt())
			testName := args
			if Subset(testName, testNames) {
				runTestsRemote(cmd, testID, "test", args, count)
			} else {
				err := fmt.Errorf("The test ID=%s:Name=%s does not exist", testID, testName)
				exitError(err)
			}
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster on which to run the test")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().IntP("count", "n", 0, "run tests n times")
	cmd.Flags().IntP("timeout", "t", 60*10, "test timeout in seconds")
	return cmd
}

func getRunTestSuiteCommand(registry *runner.TestRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "test-suite [suite]",
		Aliases: []string{"suite"},
		Short:   "Run a suite of tests on Kubernetes",
		Run: func(cmd *cobra.Command, args []string) {
			count, _ := cmd.Flags().GetInt("count")
			testSuiteID := fmt.Sprintf("test-%d", newUUIDInt())
			testSuiteNames := test.Registry.GetTestSuiteNames()
			testSuiteName := args
			if Subset(testSuiteName, testSuiteNames) {
				runTestsRemote(cmd, testSuiteID, "test-suite", args, count)
			} else {
				err := fmt.Errorf("The test suite ID=%s:Name=%s does not exist", testSuiteID, testSuiteName)
				exitError(err)
			}
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster on which to run the test")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().IntP("count", "n", 0, "run tests n times")
	cmd.Flags().IntP("timeout", "t", 60*10, "test timeout in seconds")

	return cmd
}

func getRunBenchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bench [tests]",
		Aliases: []string{"bench", "benchmarks"},
		Short:   "Run benchmarks on Kubernetes",
		Run: func(cmd *cobra.Command, args []string) {
			count, _ := cmd.Flags().GetInt("count")
			testBenchID := fmt.Sprintf("bench-%d", newUUIDInt())
			testBenchNames := test.Registry.GetBenchmarkNames()
			testBenchName := args
			if Subset(testBenchName, testBenchNames) {
				runTestsRemote(cmd, testBenchID, "bench", args, count)
			} else {
				err := fmt.Errorf("The benchmark ID=%s:Name=%s does not exist", testBenchID, testBenchName)
				exitError(err)
			}
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster on which to run the test")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().IntP("count", "n", 0, "the number of iterations to run")
	cmd.Flags().IntP("timeout", "t", 60*10, "test timeout in seconds")
	return cmd
}

func getRunBenchSuiteCommand(registry *runner.TestRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bench-suite [suite]",
		Aliases: []string{"benchmark-suite", "benchmark-suites", "bench-suites"},
		Short:   "Run a suite of benchmarks on Kubernetes",
		Run: func(cmd *cobra.Command, args []string) {
			count, _ := cmd.Flags().GetInt("count")
			testBenchSuiteID := fmt.Sprintf("bench-%d", newUUIDInt())
			testBenchSuiteNames := test.Registry.GetBenchSuiteNames()
			testBenchSutiteName := args
			if Subset(testBenchSutiteName, testBenchSuiteNames) {
				runTestsRemote(cmd, testBenchSuiteID, "bench-suite", args, count)
			} else {
				err := fmt.Errorf("The benchmark suite ID=%s:Name=%s does not exist", testBenchSuiteID, testBenchSutiteName)
				exitError(err)
			}
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster on which to run the test")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().IntP("count", "n", 0, "the number of iterations to run")
	cmd.Flags().IntP("timeout", "t", 60*10, "test timeout in seconds")
	return cmd
}

func runTestsRemote(cmd *cobra.Command, testID string, commandType string, tests []string, count int) {
	// Get the onit controller
	controller, err := k8s.NewController()
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
	if count > 0 {
		tests = append(tests, fmt.Sprintf("-n=%d", count))
	}

	message, code, status := cluster.RunTests(testID, append([]string{commandType}, tests...), time.Duration(timeout)*time.Second)
	if status.Failed() {
		exitStatus(status)
	} else {
		fmt.Println(message)
		os.Exit(code)
	}

}
