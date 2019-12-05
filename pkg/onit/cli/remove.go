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
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/onit/env"
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
		Use:   "network [args]",
		Short: "Remove a stratum network from the cluster",
		Args:  cobra.NoArgs,
		RunE:  runInCluster(runRemoveNetworkCommand),
	}
	cmd.Flags().StringP("name", "n", "", "the name of the network to remove")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_networks"},
	}
	return cmd
}

func runRemoveNetworkCommand(cmd *cobra.Command, _ []string) error {
	setupCommand(cmd)

	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	network := env.Network(getName(cmd))
	if network == nil {
		return fmt.Errorf("unknown network: %s", getName(cmd))
	}
	return network.Remove()
}

// getRemoveSimulatorCommand returns a cobra command for tearing down a device simulator
func getRemoveSimulatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulator [args]",
		Short: "Remove a device simulator from the cluster",
		Args:  cobra.NoArgs,
		RunE:  runInCluster(runRemoveSimulatorCommand),
	}
	cmd.Flags().StringP("name", "n", "", "the name of the simulator to remove")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_simulators"},
	}
	return cmd
}

func runRemoveSimulatorCommand(cmd *cobra.Command, args []string) error {
	setupCommand(cmd)

	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	simulator := env.Simulator(getName(cmd))
	if simulator == nil {
		return fmt.Errorf("unknown device: %s", getName(cmd))
	}
	return simulator.Remove()
}

// getRemoveAppCommand returns a cobra command for tearing down an app
func getRemoveAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app [args]",
		Short: "Remove an app from the cluster",
		Args:  cobra.NoArgs,
		RunE:  runInCluster(runRemoveAppCommand),
	}
	cmd.Flags().StringP("name", "n", "", "the name of the app to remove")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_apps"},
	}
	return cmd
}

func runRemoveAppCommand(cmd *cobra.Command, _ []string) error {
	setupCommand(cmd)

	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	app := env.App(getName(cmd))
	if app == nil {
		return fmt.Errorf("unknown application: %s", getName(cmd))
	}
	return app.Remove()
}
