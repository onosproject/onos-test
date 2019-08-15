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
	"sync"

	"github.com/onosproject/onos-test/pkg/onit/k8s"
	"github.com/spf13/cobra"
)

var (
	debugExample = `
		# Debug a node in the cluster
		onit debug <name of a node>

        # Debug all nodes in the cluster
        onit debug`
)

// getDebugCommand returns a cobra "debug" command to open a debugger port to the given resource
func getDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "debug <resource>",
		Short:   "Open a debugger port to the given resource",
		Example: debugExample,
		Args:    cobra.MaximumNArgs(1),
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

			nodes, err := cluster.GetNodes()
			if err != nil {
				exitError(err)
			}

			port, _ := cmd.Flags().GetInt("port")

			if len(args) == 0 {
				var wg sync.WaitGroup
				n := len(nodes)
				wg.Add(n)

				asyncErrors := make(chan error, n)
				freePorts, err := k8s.GetFreePorts(n)
				if err != nil {
					exitError(err)
				}

				for index := range nodes {
					go func(node k8s.NodeInfo, port int) {
						fmt.Println("Start Debugger for:", node.ID)
						err = cluster.PortForward(node.ID, port, 40000)
						asyncErrors <- err
						wg.Done()
					}(nodes[index], freePorts[index])

				}

				go func() {
					wg.Wait()
					close(asyncErrors)
				}()

				for err = range asyncErrors {
					if err != nil {
						exitError(err)
					}
				}

			} else {
				if err := cluster.PortForward(args[0], port, 40000); err != nil {
					exitError(err)
				} else {
					fmt.Println(port)
				}

			}
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster for which to debug nodes")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().IntP("port", "p", k8s.DebugPort, "the local port to forward to the given resource")
	return cmd
}
