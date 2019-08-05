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

	"github.com/onosproject/onos-test/pkg/onit"
	"github.com/spf13/cobra"
)

var (
	debugExample = ` 
		# Debug a node in the cluster
		onit debug <name of a node>`
)

// getDebugCommand returns a cobra "debug" command to open a debugger port to the given resource
func getDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "debug <resource>",
		Short:   "Open a debugger port to the given resource",
		Example: debugExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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

			port, _ := cmd.Flags().GetInt("port")
			if err := cluster.PortForward(args[0], port, 40000); err != nil {
				exitError(err)
			} else {
				fmt.Println(port)
			}
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster for which to load the history")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().IntP("port", "p", 40000, "the local port to forward to the given resource")
	return cmd
}
