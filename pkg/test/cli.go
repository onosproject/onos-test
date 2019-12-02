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

package test

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"math/rand"
	"time"
)

// GetCommand returns the test command
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubetest",
		Short: "Start and manage Kubernetes tests",
	}
	cmd.AddCommand(getRunCommand())
	cmd.AddCommand(getTestCommand())
	cmd.AddCommand(getBenchCommand())
	return cmd
}

// getTestCommand returns the test command
func getTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run a test",
		RunE:  runTestCommand,
	}
	cmd.Flags().StringP("image", "i", "", "the test image to run")
	cmd.Flags().StringP("suite", "s", "", "the name of a suite to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the image pull policy to use")
	cmd.Flags().Duration("timeout", 0, "the test timeout")
	return cmd
}

// runTestCommand runs the test command
func runTestCommand(cmd *cobra.Command, _ []string) error {
	image, _ := cmd.Flags().GetString("image")
	pullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	suite, _ := cmd.Flags().GetString("suite")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	config := &TestConfig{
		JobConfig: &JobConfig{
			JobID:      random.NewPetName(2),
			Type:       TestTypeTest,
			Image:      image,
			Timeout:    timeout,
			PullPolicy: corev1.PullPolicy(pullPolicy),
		},
		Suite: suite,
	}

	runner, err := NewTestRunner(config)
	if err != nil {
		return err
	}
	return runner.Run()
}

// getBenchCommand returns the bench command
func getBenchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bench",
		Short: "Run a benchmark",
		RunE:  runBenchCommand,
	}
	cmd.Flags().StringP("image", "i", "", "the test image to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the image pull policy to use")
	cmd.Flags().StringP("suite", "s", "", "the benchmark suite to run")
	cmd.Flags().StringP("benchmark", "b", "", "the name of the benchmark to run")
	cmd.Flags().IntP("clients", "p", 1, "the number of clients to run")
	cmd.Flags().Int("parallel", 1, "the number of concurrent goroutines per client")
	cmd.Flags().IntP("requests", "n", 1, "the number of requests to run")
	cmd.Flags().StringToStringP("args", "a", map[string]string{}, "a mapping of named benchmark arguments")
	cmd.Flags().Duration("timeout", 0, "the test timeout")
	return cmd
}

// runBenchCommand runs the bench command
func runBenchCommand(cmd *cobra.Command, _ []string) error {
	image, _ := cmd.Flags().GetString("image")
	pullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	suite, _ := cmd.Flags().GetString("suite")
	benchmark, _ := cmd.Flags().GetString("benchmark")
	clients, _ := cmd.Flags().GetInt("clients")
	parallelism, _ := cmd.Flags().GetInt("parallel")
	requests, _ := cmd.Flags().GetInt("requests")
	args, _ := cmd.Flags().GetStringToString("args")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	config := &BenchmarkConfig{
		JobConfig: &JobConfig{
			JobID:      random.NewPetName(2),
			Type:       TestTypeTest,
			Image:      image,
			Timeout:    timeout,
			PullPolicy: corev1.PullPolicy(pullPolicy),
		},
		Suite:       suite,
		Benchmark:   benchmark,
		Clients:     clients,
		Parallelism: parallelism,
		Requests:    requests,
		Args:        args,
	}

	runner, err := NewTestRunner(config)
	if err != nil {
		return err
	}
	return runner.Run()
}

// getRunCommand returns the run command
func getRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a test",
		RunE:  runRunCommand,
	}
	cmd.Flags().StringP("type", "t", "", "the type of test to run")
	cmd.Flags().StringP("image", "i", "", "the test image to run")
	cmd.Flags().StringP("suite", "s", "", "the name of a suite to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the image pull policy to use")
	cmd.Flags().Duration("timeout", 0, "the test timeout")
	return cmd
}

// runRunCommand runs the run command
func runRunCommand(cmd *cobra.Command, args []string) error {
	typeName, _ := cmd.Flags().GetString("type")
	switch typeName {
	case string(TestTypeTest):
		return runTestCommand(cmd, []string{})
	case string(TestTypeBenchmark):
		return runBenchCommand(cmd, []string{})
	default:
		return fmt.Errorf("unknown test type %s", typeName)
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
