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
	"github.com/onosproject/onos-test/pkg/new/kube"
	"github.com/onosproject/onos-test/pkg/new/onit/env"
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
		RunE:  runRemoveNetworkCommand,
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to which to add the network")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	_ = cmd.MarkFlagRequired("cluster")
	return cmd
}

func runRemoveNetworkCommand(cmd *cobra.Command, args []string) error {
	runCommand(cmd)

	networkID := args[0]
	cluster, _ := cmd.Flags().GetString("cluster")

	kubeAPI := kube.GetAPI(cluster)
	env := env.New(kubeAPI)
	network := env.Simulator(networkID)
	if network == nil {
		return fmt.Errorf("unknown network: %s", networkID)
	}
	return network.Remove()
}

// getRemoveSimulatorCommand returns a cobra command for tearing down a device simulator
func getRemoveSimulatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulator <name>",
		Short: "Remove a device simulator from the cluster",
		Args:  cobra.ExactArgs(1),
		RunE:  runRemoveSimulatorCommand,
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to which to add the simulator")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	_ = cmd.MarkFlagRequired("cluster")
	return cmd
}

func runRemoveSimulatorCommand(cmd *cobra.Command, args []string) error {
	runCommand(cmd)

	deviceID := args[0]
	cluster, _ := cmd.Flags().GetString("cluster")

	kubeAPI := kube.GetAPI(cluster)
	env := env.New(kubeAPI)
	simulator := env.Simulator(deviceID)
	if simulator == nil {
		return fmt.Errorf("unknown device: %s", deviceID)
	}
	return simulator.Remove()
}

// getRemoveAppCommand returns a cobra command for tearing down an app
func getRemoveAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app <name>",
		Short: "Remove an app from the cluster",
		Args:  cobra.ExactArgs(1),
		RunE:  runRemoveAppCommand,
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to which to remove the app")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	_ = cmd.MarkFlagRequired("cluster")
	return cmd
}

func runRemoveAppCommand(cmd *cobra.Command, args []string) error {
	runCommand(cmd)

	appID := args[0]
	cluster, _ := cmd.Flags().GetString("cluster")

	kubeAPI := kube.GetAPI(cluster)
	env := env.New(kubeAPI)
	app := env.App(appID)
	if app == nil {
		return fmt.Errorf("unknown application: %s", appID)
	}
	return app.Remove()
}