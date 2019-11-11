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
	"github.com/onosproject/onos-test/pkg/new/kubetest"
	"github.com/onosproject/onos-test/pkg/new/util/random"
	corev1 "k8s.io/api/core/v1"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	runExample = `
		# Run a single test on the cluster
		onit run test <name of a test>

		# Run a benchmark on the cluster
		onit run bench <name of a benchmark>`
)

// getRunCommand returns a cobra run command to run integration tests
func getRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run {test,bench}",
		Short:   "Run integration tests",
		Example: runExample,
	}
	cmd.AddCommand(getRunTestCommand())
	cmd.AddCommand(getRunBenchCommand())
	return cmd
}

func getRunTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run tests on Kubernetes",
		RunE:  runRunTestCommand,
	}
	cmd.Flags().StringP("cluster", "c", "", "the cluster on which to run the test")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().StringP("image", "i", "", "the test image to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringP("suite", "s", "", "the test suite to run")
	cmd.Flags().StringP("test", "t", "", "the name of the test method to run")
	cmd.Flags().Duration("timeout", 10*time.Minute, "test timeout")
	return cmd
}

// runRunTestCommand runs the run test command
func runRunTestCommand(cmd *cobra.Command, _ []string) error {
	return runTest(cmd, kubetest.TestTypeTest)
}

func getRunBenchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "benchmark",
		Aliases: []string{"bench"},
		Short:   "Run benchmarks on Kubernetes",
		RunE:    runRunBenchCommand,
	}
	cmd.Flags().StringP("cluster", "c", "", "the cluster on which to run the benchmark")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().StringP("image", "i", "", "the benchmark image to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringP("suite", "s", "", "the benchmark suite to run")
	cmd.Flags().StringP("benchmark", "t", "", "the name of the benchmark method to run")
	cmd.Flags().Duration("timeout", 10*time.Minute, "benchmark timeout")
	return cmd
}

// runRunBenchCommand runs the run benchmark command
func runRunBenchCommand(cmd *cobra.Command, _ []string) error {
	return runTest(cmd, kubetest.TestTypeBenchmark)
}

func runTest(cmd *cobra.Command, testType kubetest.TestType) error {
	runCommand(cmd)

	clusterID, _ := cmd.Flags().GetString("cluster")
	image, _ := cmd.Flags().GetString("image")
	suite, _ := cmd.Flags().GetString("suite")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	test := &kubetest.TestConfig{
		TestID:     random.NewPetName(2),
		Type:       testType,
		Image:      image,
		Suite:      suite,
		Timeout:    timeout,
		PullPolicy: pullPolicy,
	}

	// If the cluster ID was not specified, create a new cluster to run the test
	// Otherwise, deploy the test in the existing cluster
	if clusterID == "" {
		runner, err := kubetest.NewTestRunner(test)
		if err != nil {
			return err
		}
		return runner.Run()
	}

	cluster := kubetest.NewTestCluster(clusterID)
	if err := cluster.StartTest(test); err != nil {
		return err
	}
	if err := cluster.AwaitTestComplete(test); err != nil {
		return err
	}
	_, code, err := cluster.GetTestResult(test)
	if err != nil {
		return err
	}
	os.Exit(code)
	return nil
}
