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
	"path/filepath"

	"github.com/onosproject/onos-test/pkg/onit/k8s"

	"github.com/onosproject/onos-test/pkg/onit/console"
	"github.com/spf13/cobra"
)

var (
	fetchExample = `
		# Download logs from all nodes
		onit fetch logs 

		# Download logs from a node
		onit fetch logs <name of the node>`
)

// getFetchCommand returns a cobra "download" command for downloading resources from a test cluster
func getFetchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fetch",
		Short:   "Fetch resources from the cluster",
		Example: fetchExample,
	}
	cmd.AddCommand(getFetchLogsCommand())
	return cmd
}

// getFetchLogsCommand returns a cobra command for downloading the logs from a node
func getFetchLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [node]",
		Short: "Download logs from a node",
		Args:  cobra.MaximumNArgs(1),
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

			options := parseLogOptions(cmd)

			destination, _ := cmd.Flags().GetString("destination")
			if len(args) > 0 {
				resourceID := args[0]
				resources, err := cluster.GetResources(resourceID)
				if err != nil {
					exitError(err)
				}

				var status console.ErrorStatus
				for _, resource := range resources {
					path := filepath.Join(destination, fmt.Sprintf("%s.log", resource))
					status = cluster.DownloadLogs(resource, path, options)
				}

				if status.Failed() {
					exitStatus(status)
				}
			} else {
				nodes, err := cluster.GetNodes()
				if err != nil {
					exitError(err)
				}

				var status console.ErrorStatus
				for _, node := range nodes {
					path := filepath.Join(destination, fmt.Sprintf("%s.log", node.ID))
					status = cluster.DownloadLogs(node.ID, path, options)
				}

				if status.Failed() {
					exitStatus(status)
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
	cmd.Flags().StringP("destination", "d", ".", "the destination to which to write the logs")
	return cmd
}
