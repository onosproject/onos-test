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
	"errors"

	"github.com/onosproject/onos-test/pkg/onit/k8s"

	"github.com/spf13/cobra"
)

var (
	removeExample = `
		# Remove a simulator with a given name
		onit remove simulator <simulator-name>

		# Remove a network with a given name
		onit remove network <network-name>
	
		# Remove an app
		onit remove app <app-name>`
)

// getRemoveCommand returns a cobra "remove" command for removing resources from the cluster
func getRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove {simulator} [args]",
		Short:   "Remove resources from the cluster",
		Example: removeExample,
	}
	cmd.AddCommand(getRemoveSimulatorCommand())
	cmd.AddCommand(getRemoveNetworkCommand())
	cmd.AddCommand(getRemoveAppCommand())
	return cmd
}

// getRemoveNetworkCommand returns a cobra command for tearing down a stratum network
func getRemoveNetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network [name]",
		Short: "Remove a stratum network from the cluster",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

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

			networks, err := cluster.GetNetworks()
			if err != nil {
				exitError(err)
			}
			if !Contains(networks, name) {
				exitError(errors.New("The given network name does not exist"))
			}

			// Remove the network from the cluster
			if status := cluster.RemoveNetwork(name); status.Failed() {
				exitStatus(status)
			}
		},
	}

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to which to add the network")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}

// getRemoveSimulatorCommand returns a cobra command for tearing down a device simulator
func getRemoveSimulatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulator <name>",
		Short: "Remove a device simulator from the cluster",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

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

			simulators, err := cluster.GetSimulators()
			if err != nil {
				exitError(err)
			}

			if !Contains(simulators, name) {
				exitError(errors.New("The given simulator name does not exist"))
			}

			// Remove the simulator from the cluster
			if status := cluster.RemoveSimulator(name); status.Failed() {
				exitStatus(status)
			}
		},
	}

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to which to add the simulator")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}

// getRemoveAppCommand returns a cobra command for tearing down an app
func getRemoveAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app <name>",
		Short: "Remove an app from the cluster",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

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

			apps, err := cluster.GetApps()
			if err != nil {
				exitError(err)
			}

			if !Contains(apps, name) {
				exitError(errors.New("the given app name does not exist"))
			}

			// Remove the app from the cluster
			if status := cluster.RemoveApp(name); status.Failed() {
				exitStatus(status)
			}
		},
	}

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to which to remove the app")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}
