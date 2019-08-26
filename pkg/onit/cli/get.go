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
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/onosproject/onos-test/pkg/onit/k8s"

	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var (
	getExample = ` 
		# Get current configured cluster
		onit get cluster

		# Get a list of clusters
		onit get clusters
            
		# Get a list of nodes
		onit get nodes

		# Get a list of simulators
		onit get simulators

		# Get a list of networks
		onit get networks

		# Get a list of partitions
		onit get partitions

		# Get a list of nodes in a partition
		onit get partition <partition-id>
            
		# Get a list of integration tests
		onit get tests

		# Get a list of integration testing suites
		onit get test-suites

		# Get a list of benchmarks
		onit get benchmarks

		# Get a list of benchmark suites
		onit get bench-suites
            
		# Get the logs for a test resource
		onit get logs <name of resource>
            
		# Get the history of test runs
		onit get history

		# Get the list of installed apps
		onit get apps`
)

// getGetCommand returns a cobra "get" command to read test configurations
func getGetCommand(registry *runner.TestRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get {cluster,clusters,networks,simulators,device-presets,store-presets,tests,suites}",
		Short:   "Get test configurations",
		Example: getExample,
	}
	cmd.AddCommand(getGetClusterCommand())
	cmd.AddCommand(getGetNodesCommand())
	cmd.AddCommand(getGetPartitionsCommand())
	cmd.AddCommand(getGetPartitionCommand())
	cmd.AddCommand(getGetSimulatorsCommand())
	cmd.AddCommand(getGetNetworksCommand())
	cmd.AddCommand(getGetClustersCommand())
	cmd.AddCommand(getGetDevicePresetsCommand())
	cmd.AddCommand(getGetStorePresetsCommand())
	cmd.AddCommand(getGetTestsCommand(registry))
	cmd.AddCommand(getGetTestSuitesCommand(registry))
	cmd.AddCommand(getGetBenchmarksCommand(registry))
	cmd.AddCommand(getGetBenchmarkSuitesCommand(registry))
	cmd.AddCommand(getGetHistoryCommand())
	cmd.AddCommand(getGetLogsCommand())
	cmd.AddCommand(getGetAppsCommand())
	return cmd
}

// getGetClusterCommand returns a cobra command to get the current cluster context
func getGetClusterCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cluster",
		Short: "Get the currently configured cluster",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getDefaultCluster())
		},
	}
}

// getGetNetworksCommand returns a cobra command to get the list of networks deployed in the current cluster context
func getGetNetworksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "networks",
		Short: "Get the currently configured cluster's networks",
		Run: func(cmd *cobra.Command, args []string) {
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

			// Get the list of networks and output
			networks, err := cluster.GetNetworks()
			if err != nil {
				exitError(err)
			} else {
				for _, name := range networks {
					fmt.Println(name)
				}
			}
		},
	}

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to query")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}

// getGetSimulatorsCommand returns a cobra command to get the list of simulators deployed in the current cluster context
func getGetSimulatorsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulators",
		Short: "Get the currently configured cluster's simulators",
		Run: func(cmd *cobra.Command, args []string) {
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

			// Get the list of simulators and output
			simulators, err := cluster.GetSimulators()
			if err != nil {
				exitError(err)
			} else {
				for _, name := range simulators {
					fmt.Println(name)
				}
			}
		},
	}

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to query")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}

// getGetClustersCommand returns a cobra command to get a list of available test clusters
func getGetClustersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clusters",
		Short: "Get a list of all deployed clusters",
		Run: func(cmd *cobra.Command, args []string) {
			// Get the onit controller
			controller, err := k8s.NewController()
			if err != nil {
				exitError(err)
			}

			// Get the list of clusters and output
			clusters, err := controller.GetClusters()
			if err != nil {
				exitError(err)
			} else {
				noHeaders, _ := cmd.Flags().GetBool("no-headers")
				printClusters(clusters, !noHeaders)
			}
		},
	}
	cmd.Flags().Bool("no-headers", false, "whether to print column headers")
	return cmd
}

func printClusters(clusters map[string]*k8s.ClusterConfig, includeHeaders bool) {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	if includeHeaders {
		fmt.Fprintln(writer, "ID\tSIZE\tPARTITIONS")
	}
	for id, config := range clusters {
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%d\t%d\t%d", id, config.ConfigNodes, config.TopoNodes, config.Partitions))
	}
	writer.Flush()
}

// getGetDevicePresetsCommand returns a cobra command to get a list of available device simulator configurations
func getGetDevicePresetsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "device-presets",
		Short: "Get a list of device configurations",
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range getSimulatorPresets() {
				fmt.Println(name)
			}
		},
	}
}

// getGetStorePresetsCommand returns a cobra command to get a list of available store configurations
func getGetStorePresetsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "store-presets",
		Short: "Get a list of store configurations",
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range getStorePresets() {
				fmt.Println(name)
			}
		},
	}
}

// getGetPartitionsCommand returns a cobra command to get a list of Raft partitions in the cluster
func getGetPartitionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "partitions",
		Short: "Get a list of partitions in the cluster",
		Run: func(cmd *cobra.Command, args []string) {
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

			// Get the list of partitions and output
			partitions, err := cluster.GetPartitions()
			if err != nil {
				exitError(err)
			} else {
				printPartitions(partitions)
			}
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to query")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}

func printPartitions(partitions []k8s.PartitionInfo) {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	fmt.Fprintln(writer, "ID\tGROUP\tNODES")
	for _, partition := range partitions {
		fmt.Fprintln(writer, fmt.Sprintf("%d\t%s\t%s", partition.Partition, partition.Group, strings.Join(partition.Nodes, ",")))
	}
	writer.Flush()
}

// getGetAppsCommand returns a cobra command to get the list of apps deployed in the current cluster context
func getGetAppsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: "Get the currently configured cluster's apps",
		Run: func(cmd *cobra.Command, args []string) {
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

			// Get the list of simulators and output
			apps, err := cluster.GetApps()
			if err != nil {
				exitError(err)
			} else {
				for _, name := range apps {
					fmt.Println(name)
				}
			}
		},
	}

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to query")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}

// getGetPartitionCommand returns a cobra command to get the nodes in a partition
func getGetPartitionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "partition <partition>",
		Short: "Get a list of nodes in a partition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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

			// Parse the partition argument
			partition, err := strconv.ParseInt(args[0], 0, 32)
			if err != nil {
				exitError(err)
			}

			// Get the list of nodes and output
			nodes, err := cluster.GetPartitionNodes(int(partition))
			if err != nil {
				exitError(err)
			} else {
				printNodes(nodes, true)
			}
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to query")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}

// getGetNodesCommand returns a cobra command to get a list of onos nodes in the cluster
func getGetNodesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Get a list of nodes in the cluster",
		Run: func(cmd *cobra.Command, args []string) {
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

			nodeType, err := cmd.Flags().GetString("type")
			if err != nil {
				exitError(err)
			}

			// Get the cluster controller
			cluster, err := controller.GetCluster(clusterID)
			if err != nil {
				exitError(err)
			}

			// Get the list of nodes and output
			if strings.Compare(nodeType, string(k8s.OnosAll)) == 0 {
				nodes, err := cluster.GetNodes()
				if err != nil {
					exitError(err)
				} else {
					noHeaders, _ := cmd.Flags().GetBool("no-headers")
					printNodes(nodes, !noHeaders)
				}
			} else if strings.Compare(nodeType, string(k8s.OnosConfig)) == 0 {
				nodes, err := cluster.GetOnosConfigNodes()
				if err != nil {
					exitError(err)
				} else {
					noHeaders, _ := cmd.Flags().GetBool("no-headers")
					printNodes(nodes, !noHeaders)
				}

			} else if strings.Compare(nodeType, string(k8s.OnosTopo)) == 0 {
				nodes, err := cluster.GetOnosTopoNodes()
				if err != nil {
					exitError(err)
				} else {
					noHeaders, _ := cmd.Flags().GetBool("no-headers")
					printNodes(nodes, !noHeaders)
				}

			} else if strings.Compare(nodeType, string(k8s.OnosCli)) == 0 {
				nodes, err := cluster.GetOnosCliNodes()
				if err != nil {
					exitError(err)
				} else {
					noHeaders, _ := cmd.Flags().GetBool("no-headers")
					printNodes(nodes, !noHeaders)
				}
			} else if strings.Compare(nodeType, string(k8s.OnosGui)) == 0 {
				nodes, err := cluster.GetOnosGuiNodes()
				if err != nil {
					exitError(err)
				} else {
					noHeaders, _ := cmd.Flags().GetBool("no-headers")
					printNodes(nodes, !noHeaders)
				}

			}
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to query")
	cmd.Flags().StringP("type", "t", "all", "To get list of nodes based on their types")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().Bool("no-headers", false, "whether to print column headers")
	return cmd
}

func printNodes(nodes []k8s.NodeInfo, includeHeaders bool) {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	if includeHeaders {
		fmt.Fprintln(writer, "ID\tTYPE\tSTATUS")
	}
	for _, node := range nodes {
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s", node.ID, node.Type, node.Status))
	}
	writer.Flush()
}

// getGetTestsCommand returns a cobra command to get a list of available tests
func getGetTestsCommand(registry *runner.TestRegistry) *cobra.Command {
	return &cobra.Command{
		Use:   "tests",
		Short: "Get a list of integration tests",
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range registry.GetTestNames() {
				fmt.Println(name)
			}
		},
	}
}

// getGetTestsCommand returns a cobra command to get a list of available tests
func getGetTestSuitesCommand(registry *runner.TestRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "test-suites",
		Aliases: []string{"suites"},
		Short:   "Get a list of integration testing suites",
		Run: func(cmd *cobra.Command, args []string) {
			noHeaders, _ := cmd.Flags().GetBool("no-headers")
			printTestSuites(registry, !noHeaders)
		},
	}

	cmd.Flags().Bool("no-headers", false, "whether to print column headers")
	return cmd
}

// PrintTestSuites prints test suites in a table
func printTestSuites(registry *runner.TestRegistry, includeHeaders bool) {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	if includeHeaders {
		fmt.Fprintln(writer, "SUITE\tTESTS")
	}
	for name, suite := range registry.TestSuites {
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%s", name, strings.Join(suite.GetTestNames(), ", ")))
	}
	writer.Flush()
}

// getGetBenchmarksCommand returns a cobra command to get a list of available tests
func getGetBenchmarksCommand(registry *runner.TestRegistry) *cobra.Command {
	return &cobra.Command{
		Use:     "benchmarks",
		Aliases: []string{"bench", "benchmark"},
		Short:   "Get a list of benchmarks",
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range registry.GetBenchmarkNames() {
				fmt.Println(name)
			}
		},
	}
}

// getGetBenchmarkSuitesCommand returns a cobra command to get a list of available tests
func getGetBenchmarkSuitesCommand(registry *runner.TestRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bench-suites",
		Aliases: []string{"benchmark-suites"},
		Short:   "Get a list of benchmark suites",
		Run: func(cmd *cobra.Command, args []string) {
			noHeaders, _ := cmd.Flags().GetBool("no-headers")
			printBenchSuites(registry, !noHeaders)
		},
	}

	cmd.Flags().Bool("no-headers", false, "whether to print column headers")
	return cmd
}

// printBenchSuites prints benchmark suites in a table
func printBenchSuites(registry *runner.TestRegistry, includeHeaders bool) {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	if includeHeaders {
		fmt.Fprintln(writer, "SUITE\tBENCHMARKS")
	}
	for name, suite := range registry.BenchSuites {
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%s", name, strings.Join(suite.GetBenchNames(), ", ")))
	}
	writer.Flush()
}

// getGetHistoryCommand returns a cobra command to get the history of tests
func getGetHistoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Get the history of test runs",
		Run: func(cmd *cobra.Command, args []string) {
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

			// Get the history of test runs for the cluster
			records, err := cluster.GetHistory()
			if err != nil {
				exitError(err)
			}

			printHistory(records)
		},
	}

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster for which to load the history")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}

// printHistory prints a test history in table format
func printHistory(records []k8s.TestRecord) {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	fmt.Fprintln(writer, "ID\tTESTS\tSTATUS\tEXIT CODE\tMESSAGE")
	for _, record := range records {
		var args string
		if len(record.Args) > 0 {
			args = strings.Join(record.Args, ",")
		} else {
			args = "*"
		}
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s\t%d\t%s", record.TestID, args, record.Status, record.ExitCode, record.Message))
	}
	writer.Flush()
}

// getGetLogsCommand returns a cobra command to output the logs for a specific resource
func getGetLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs <id> [options]",
		Short: "Get the logs for a test resource",
		Long: `Outputs the complete logs for any test resource.
To output the logs from an onos-config node, get the node ID via 'onit get nodes'
To output the logs from a test, get the test ID from the test run or from 'onit get history'`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stream, _ := cmd.Flags().GetBool("stream")

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

			// If streaming is enabled, stream the logs to stdout. Otherwise, get the logs for all resources and print them.
			if stream {
				reader, err := cluster.StreamLogs(args[0])
				if err != nil {
					exitError(err)
				}
				defer reader.Close()
				if err = printStream(reader); err != nil {
					exitError(err)
				}
			} else {
				resources, err := cluster.GetResources(args[0])
				if err != nil {
					exitError(err)
				}

				// Iterate through resources and get/print logs
				numResources := len(resources)

				if err != nil {
					exitError(err)
				}

				options := parseLogOptions(cmd)

				for i, resource := range resources {
					logs, err := cluster.GetLogs(resource, options)
					if err != nil {
						exitError(err)
					}
					os.Stdout.Write(logs)
					if i+1 < numResources {
						fmt.Println("----")
					}
				}
			}
		},
	}

	cmd.Flags().DurationP("since", "", -1, "Only return logs newer than a relative "+
		"duration like 5s, 2m, or 3h. Defaults to all logs. Only one of since-time / since may be used")
	cmd.Flags().Int64P("tail", "t", -1, "If set, the number of bytes to read from the "+
		"server before terminating the log output. This may not display a complete final line of logging, and may return "+
		"slightly more or slightly less than the specified limit.")

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster for which to load the history")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().BoolP("stream", "s", false, "stream logs to stdout")
	return cmd
}

func parseLogOptions(cmd *cobra.Command) corev1.PodLogOptions {
	// Parse log options from CLI
	options := corev1.PodLogOptions{}
	since, err := cmd.Flags().GetDuration("since")
	sinceSeconds := int64(since / time.Second)
	if sinceSeconds > 0 {
		options.SinceSeconds = &sinceSeconds
	}
	if err != nil {
		exitError(err)
	}
	tail, err := cmd.Flags().GetInt64("tail")
	if tail > 0 && err != nil {
		options.TailLines = &tail
	}
	return options
}

// printStream prints a stream to stdout from the given reader
func printStream(reader io.Reader) error {
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fmt.Print(string(buf[:n]))
	}
	return nil
}
