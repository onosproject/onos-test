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
	"github.com/onosproject/onos-test/pkg/onit/setup"

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

			// Get the cluster ID
			clusterID, err := cmd.Flags().GetString("cluster")
			if err != nil {
				exitError(err)
			}

			testSetupBuilder := setup.New()
			testSetupBuilder.SetClusterID(clusterID).SetNetworkName(name)
			testSetup := testSetupBuilder.Build()
			testSetup.RemoveNetwork()

		},
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to which to add the network")
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

			// Get the cluster ID
			clusterID, err := cmd.Flags().GetString("cluster")
			if err != nil {
				exitError(err)
			}

			testSetupBuilder := setup.New()
			testSetupBuilder.SetClusterID(clusterID).SetSimulatorName(name)
			testSetup := testSetupBuilder.Build()
			testSetup.RemoveSimulator()

		},
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to which to add the simulator")
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

			// Get the cluster ID
			clusterID, err := cmd.Flags().GetString("cluster")
			if err != nil {
				exitError(err)
			}

			testSetupBuilder := setup.New()
			testSetupBuilder.SetClusterID(clusterID).SetAppName(name)
			testSetup := testSetupBuilder.Build()
			testSetup.RemoveApp()

		},
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to which to remove the app")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}
