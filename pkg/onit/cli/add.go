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

	"github.com/onosproject/onos-test/pkg/onit/setup"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var (
	addExample = `
		# Add a simulator with a given name
		onit add simulator simulator-1

		# Add a network of stratum switches that emulates a linear network topology with two nodes
		onit add network stratum-linear -- --topo linear,2
	   
		# Add latest version of an application 
		onit add app onos-ztp --image onosproject/onos-ztp:latest --image-pull-policy "Always" `
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

			// Get the cluster ID
			clusterID, err := cmd.Flags().GetString("cluster")
			if err != nil {
				exitError(err)
			}

			testSetupBuilder := setup.New()
			testSetupBuilder.
				SetClusterID(clusterID).
				SetConfigName(configName).
				SetNetworkName(name)

			if len(args) > 1 {
				testSetupBuilder.SetMininetOptions(args[1:])
			}
			testSetup := testSetupBuilder.Build()
			testSetup.AddNetwork()

		},
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to which to add the simulator")
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

			testSetupBuilder := setup.New()

			// Create the simulator configuration from the configured preset
			configName, _ := cmd.Flags().GetString("preset")

			// Get the cluster ID
			clusterID, err := cmd.Flags().GetString("cluster")
			if err != nil {
				exitError(err)
			}
			testSetupBuilder.SetConfigName(configName)
			testSetupBuilder.SetSimulatorName(name)
			testSetupBuilder.SetClusterID(clusterID)
			testSetup := testSetupBuilder.Build()
			testSetup.AddSimulator()
			if err != nil {
				exitError(err)
			}

		},
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to which to add the simulator")
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
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var name string
			if len(args) == 1 {
				name = args[0]
			}

			// If the name is not set, assign a generic UUID based name.
			if name == "" {
				name = fmt.Sprintf("app-%d", newUUIDInt())
			}

			image, _ := cmd.Flags().GetString("image")
			imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")

			// Get the cluster ID
			clusterID, err := cmd.Flags().GetString("cluster")
			if err != nil {
				exitError(err)
			}

			testSetupBuilder := setup.New()
			testSetupBuilder.SetClusterID(clusterID).SetImageName(image).SetImagePullPolicy(imagePullPolicy)
			testSetupBuilder.SetAppName(name)

			testSetup := testSetupBuilder.Build()
			testSetup.AddApp()
		},
	}

	cmd.Flags().StringP("image", "i", "", "the image name")
	_ = cmd.MarkFlagRequired("image")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to which to add the app")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}
