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
	"text/tabwriter"

	"github.com/onosproject/onos-test/pkg/cluster"
	"github.com/onosproject/onos-test/pkg/kube"
	oc "github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/spf13/cobra"
)

var (
	getExample = ` 
		# Get a list of clusters
		onit get clusters

		# Get the list of installed apps
		onit get apps

		# Get a list of simulators
		onit get simulators

		# Get a list of networks
		onit get networks`
)

const benchJobHeader = "Job Name\tStatus\tType\tImage\tBenchmark Suite\tBenchmark Name"
const testJobHeader = "Job Name\tStatus\tType\tImage\tTest Suite\tTest Name"

// getGetCommand returns a cobra "get" command to read test configurations
func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get {cluster,clusters,networks,simulators, tests, benchmarks}",
		Short:   "Get test configurations",
		Example: getExample,
	}
	cmd.AddCommand(getGetClustersCommand())
	cmd.AddCommand(getGetSimulatorsCommand())
	cmd.AddCommand(getGetNetworksCommand())
	cmd.AddCommand(getGetAppsCommand())
	cmd.AddCommand(getGetTestsCommand())
	cmd.AddCommand(getGetTestCommand())
	cmd.AddCommand(getGetBenchmarksCommand())
	cmd.AddCommand(getGetBenchmarkCommand())
	return cmd
}

// getGetBenchmarkCommand returns a cobra command to get history of a specific benchmark
func getGetBenchmarkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "benchmark [name]",
		Short: "Get history of a benchmark",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runGetBenchmarkCommand,
	}
	return cmd
}

// getGetTestCommand returns a cobra command to get history of a specific test
func getGetTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test [name]",
		Short: "Get history of a test",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runGetTestCommand,
	}
	return cmd
}

func runGetBenchmarkCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return nil
	}
	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	benchmarks := env.History().GetBenchmarksMap()
	benchmark := benchmarks[args[0]]
	printBenchmark(benchmark)
	return nil
}

func runGetTestCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return nil
	}
	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	tests := env.History().GetTestsMap()
	test := tests[args[0]]
	printTest(test)
	return nil
}

// getGetTestsCommand returns a cobra command to get history of tests
func getGetTestsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tests",
		Short: "Get history of tests",
		RunE:  runGetTestsCommand,
	}
	return cmd
}

// getGetBenchmarksCommand returns a cobra command to get history of benchmarks
func getGetBenchmarksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "benchmarks",
		Short: "Get history of benchmarks",
		RunE:  runGetBenchmarksCommand,
	}
	return cmd
}

func makeJobInfoHeader(job oc.JobInfo) string {
	jobName := job.GetJobName()
	jobStatus := job.GetJobStatus()
	jobType := job.GetJobType()
	jobImage := job.GetJobImage()
	jobInfo := jobName + "\t" + jobStatus + "\t" + jobType + "\t" + jobImage + "\t"
	if job.GetJobType() == "test" {
		jobInfo = jobInfo + job.GetEnvVar()["TEST_SUITE"] + "\t" + job.GetEnvVar()["TEST_NAME"]
	}
	if job.GetJobType() == "benchmark" {
		jobInfo = jobInfo + job.GetEnvVar()["BENCHMARK_SUITE"] + "\t" + job.GetEnvVar()["BENCHMARK_NAME"]
	}
	return jobInfo
}

func printBenchmark(job oc.JobInfo) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug|tabwriter.AlignRight)
	fmt.Fprintln(w, benchJobHeader)
	fmt.Fprintln(w, makeJobInfoHeader(job))
	w.Flush()
}

func printBenchmarks(jobs []oc.JobInfo) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug|tabwriter.AlignRight)
	fmt.Fprintln(w, benchJobHeader)
	for _, job := range jobs {
		fmt.Fprintln(w, makeJobInfoHeader(job))
	}
	w.Flush()
}

func printTest(job oc.JobInfo) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug|tabwriter.AlignRight)
	fmt.Fprintln(w, testJobHeader)
	fmt.Fprintln(w, makeJobInfoHeader(job))
	w.Flush()
}

func printTests(jobs []oc.JobInfo) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug|tabwriter.AlignRight)
	_, _ = fmt.Fprintln(w, testJobHeader)
	for _, job := range jobs {
		_, _ = fmt.Fprintln(w, makeJobInfoHeader(job))
	}
	_ = w.Flush()
}

func runGetBenchmarksCommand(cmd *cobra.Command, _ []string) error {
	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}
	env := env.New(kubeAPI)
	benchmarks := env.History().ListBenchmarks()
	printBenchmarks(benchmarks)
	return nil
}

func runGetTestsCommand(cmd *cobra.Command, _ []string) error {
	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	tests := env.History().ListTests()
	printTests(tests)
	return nil
}

// getGetClustersCommand returns a cobra command to get a list of available test clusters
func getGetClustersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clusters",
		Short: "Get a list of all deployed clusters",
		RunE:  runGetClustersCommand,
	}
}

func runGetClustersCommand(cmd *cobra.Command, _ []string) error {
	clusters, err := cluster.GetClusters()
	if err != nil {
		return err
	}

	for _, cluster := range clusters {
		fmt.Println(cluster)
	}
	return nil
}

// getGetNetworksCommand returns a cobra command to get the list of networks deployed in the current cluster context
func getGetNetworksCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "networks",
		Short: "Get the currently configured cluster's networks",
		RunE:  runGetNetworksCommand,
	}
}

func runGetNetworksCommand(cmd *cobra.Command, _ []string) error {
	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	for _, network := range env.Networks().List() {
		fmt.Println(network.Name())
	}
	return nil
}

// getGetSimulatorsCommand returns a cobra command to get the list of simulators deployed in the current cluster context
func getGetSimulatorsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "simulators",
		Short: "Get the currently configured cluster's simulators",
		RunE:  runGetSimulatorsCommand,
	}
}

func runGetSimulatorsCommand(cmd *cobra.Command, _ []string) error {
	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	for _, simulator := range env.Simulators().List() {
		fmt.Println(simulator.Name())
	}
	return nil
}

// getGetAppsCommand returns a cobra command to get the list of apps deployed in the current cluster context
func getGetAppsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "apps",
		Short: "Get the currently configured cluster's apps",
		RunE:  runGetAppsCommand,
	}
}

func runGetAppsCommand(cmd *cobra.Command, _ []string) error {
	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	for _, app := range env.Apps().List() {
		fmt.Println(app.Name())
	}
	return nil
}
