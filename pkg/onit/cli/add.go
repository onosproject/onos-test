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
	addExample = `
		# Add a simulator with a given name
		onit add simulator simulator-1

		# Add a network of stratum switches that emulates a linear network topology with two nodes
		onit add network stratum-linear -- --topo linear,2`
)

// getAddCommand returns a cobra "add" command for adding resources to the cluster
func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add {simulator,network} [args]",
		Short:   "Add resources to the cluster",
		Example: addExample,
	}
	cmd.AddCommand(getAddSimulatorCommand())
	cmd.AddCommand(getAddNetworkCommand())
	cmd.AddCommand(getAddAppCommand())
	return cmd
}

// getAddNetworkCommand returns a cobra command for deploying a stratum network
func getAddNetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network [name]",
		Short: "Add a stratum network to the test cluster",
		Args:  cobra.MaximumNArgs(10),
		Run: func(cmd *cobra.Command, args []string) {
			// If a network name was not provided, generate one from a UUID.
			var name string
			if len(args) > 0 {
				name = args[0]
			} else {
				name = fmt.Sprintf("network-%d", newUUIDInt())
			}

			// Create the simulator configuration from the configured preset
			configName, _ := cmd.Flags().GetString("preset")

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

			// Create the network configuration

			config := &onit.NetworkConfig{
				Config: configName,
			}
			if len(args) > 1 {
				config.MininetOptions = args[1:]
			}

			// Update number of devices in the network configuration
			onit.ParseMininetOptions(config)

			if err != nil {
				exitError(err)
			}

			// Add the network to the cluster
			if status := cluster.AddNetwork(name, config); status.Failed() {
				exitStatus(status)
			} else {
				fmt.Println(name)
			}
		},
	}

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to which to add the simulator")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().StringP("preset", "p", "default", "simulator preset to apply")
	return cmd
}

// getAddSimulatorCommand returns a cobra command for deploying a device simulator
func getAddSimulatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulator [name]",
		Short: "Add a device simulator to the test cluster",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// If a simulator name was not provided, generate one from a UUID.
			var name string
			if len(args) > 0 {
				name = args[0]
			} else {
				name = fmt.Sprintf("device-%d", newUUIDInt())
			}

			// Create the simulator configuration from the configured preset
			configName, _ := cmd.Flags().GetString("preset")

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

			// Create the simulator configuration
			config := &onit.SimulatorConfig{
				Config: configName,
			}

			// Add the simulator to the cluster
			if status := cluster.AddSimulator(name, config); status.Failed() {
				exitStatus(status)
			} else {
				fmt.Println(name)
			}
		},
	}

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to which to add the simulator")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	cmd.Flags().StringP("preset", "p", "default", "simulator preset to apply")
	return cmd
}

// getAddSimulatorCommand returns a cobra command for deploying a device simulator
func getAddAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app image-name [name]",
		Short: "Add an app to the test cluster",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// If an app name was not provided, generate one from a UUID.
			var name string

			if len(args) > 1 {
				name = args[1]
			} else {
				name = fmt.Sprintf("app-%d", newUUIDInt())
			}

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

			// Create the app configuration
			imageName := args[0]
			config := &onit.AppConfig{
				Image: imageName,
			}

			// Add the app to the cluster
			if status := cluster.AddApp(name, config); status.Failed() {
				exitStatus(status)
			} else {
				fmt.Println(name)
			}
		},
	}

	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to which to add the app")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}
