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
	"github.com/onosproject/onos-test/pkg/test"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"os"
	"strings"
	"time"
)

func getBenchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "benchmark",
		Aliases: []string{"benchmarks", "bench"},
		Short:   "Run benchmarks on Kubernetes",
		RunE:    runBenchCommand,
	}
	cmd.Flags().StringP("image", "i", "", "the benchmark image to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringToStringP("image-override", "o", map[string]string{}, "a mapping of image overrides")
	cmd.Flags().StringP("suite", "s", "", "the benchmark suite to run")
	cmd.Flags().StringP("benchmark", "b", "", "the name of the benchmark to run")
	cmd.Flags().IntP("clients", "p", 1, "the number of clients to run")
	cmd.Flags().Int("parallel", 1, "the number of concurrent goroutines per client")
	cmd.Flags().IntP("requests", "n", 1, "the number of requests to run")
	cmd.Flags().StringToStringP("arg", "a", map[string]string{}, "a mapping of named benchmark arguments")
	cmd.Flags().Duration("timeout", 10*time.Minute, "benchmark timeout")
	cmd.Flags().Bool("no-teardown", false, "do not tear down clusters following tests")
	return cmd
}

func runBenchCommand(cmd *cobra.Command, _ []string) error {
	runCommand(cmd)

	clusterID, _ := cmd.Flags().GetString("cluster")
	image, _ := cmd.Flags().GetString("image")
	images, _ := cmd.Flags().GetStringToString("image-override")
	suite, _ := cmd.Flags().GetString("suite")
	benchmark, _ := cmd.Flags().GetString("benchmark")
	clients, _ := cmd.Flags().GetInt("clients")
	parallelism, _ := cmd.Flags().GetInt("parallel")
	requests, _ := cmd.Flags().GetInt("requests")
	args, _ := cmd.Flags().GetStringToString("arg")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	noTeardown, _ := cmd.Flags().GetBool("no-teardown")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	env := make(map[string]string)
	for name, image := range images {
		env[fmt.Sprintf("IMAGE_%s", strings.ToUpper(name))] = image
	}

	config := &test.BenchmarkConfig{
		JobConfig: &test.JobConfig{
			JobID:      random.NewPetName(2),
			Type:       test.TestTypeBenchmark,
			Image:      image,
			Env:        env,
			Timeout:    timeout,
			PullPolicy: pullPolicy,
			Teardown:   !noTeardown,
		},
		Suite:       suite,
		Benchmark:   benchmark,
		Clients:     clients,
		Parallelism: parallelism,
		Requests:    requests,
		Args:        args,
	}

	// If the cluster ID was not specified, create a new cluster to run the test
	// Otherwise, deploy the test in the existing cluster
	if clusterID == "" {
		runner, err := test.NewTestRunner(config)
		if err != nil {
			return err
		}
		return runner.Run()
	}

	cluster, err := test.NewTestCluster(clusterID)
	if err != nil {
		return err
	}
	if err := cluster.StartTest(config); err != nil {
		return err
	}
	if err := cluster.AwaitTestComplete(config); err != nil {
		return err
	}
	_, code, err := cluster.GetTestResult(config)
	if err != nil {
		return err
	}
	os.Exit(code)
	return nil
}
