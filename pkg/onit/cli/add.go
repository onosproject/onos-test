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
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/onosproject/onos-test/pkg/util/random"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var (
	addExample = `
		# Add a simulator with a given name
		onit add simulator simulator-1

		# Add a network of stratum switches that emulates a linear network topology with two nodes
		onit add network stratum-linear --topo linear,2
	   
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
		RunE:  runAddNetworkCommand,
	}

	cmd.Flags().StringP("image", "i", defaultMininetImage, "the image to deploy")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringP("cluster", "c", "", "the cluster to which to add the simulator")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	_ = cmd.MarkFlagRequired("cluster")
	cmd.Flags().StringP("topo", "t", "", "the topology to create")
	_ = cmd.MarkFlagRequired("topo")
	cmd.Flags().IntP("devices", "d", 0, "the number of devices in the topology")
	_ = cmd.MarkFlagRequired("devices")
	return cmd
}

func runAddNetworkCommand(cmd *cobra.Command, args []string) error {
	runCommand(cmd)

	var networkID string
	if len(args) > 0 {
		networkID = args[0]
	}
	if networkID == "" {
		networkID = random.NewPetName(2)
	}

	cluster, _ := cmd.Flags().GetString("cluster")
	image, _ := cmd.Flags().GetString("image")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)
	topo, _ := cmd.Flags().GetString("topo")
	devices, _ := cmd.Flags().GetInt("devices")

	kubeAPI, err := kube.GetAPI(cluster)
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	_, err = env.Networks().
		New().
		Name(networkID).
		Image(image).
		PullPolicy(pullPolicy).
		Custom(topo, devices).
		Add()
	return err
}

// getAddSimulatorCommand returns a cobra command for deploying a device simulator
func getAddSimulatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulator [name]",
		Short: "Add a device simulator to the test cluster",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runAddSimulatorCommand,
	}

	cmd.Flags().StringP("image", "i", defaultSimulatorImage, "the image to deploy")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringP("cluster", "c", "", "the cluster to which to add the simulator")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	_ = cmd.MarkFlagRequired("cluster")
	return cmd
}

func runAddSimulatorCommand(cmd *cobra.Command, args []string) error {
	runCommand(cmd)

	var deviceID string
	if len(args) > 0 {
		deviceID = args[0]
	}
	if deviceID == "" {
		deviceID = random.NewPetName(2)
	}

	cluster, _ := cmd.Flags().GetString("cluster")
	image, _ := cmd.Flags().GetString("image")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	kubeAPI, err := kube.GetAPI(cluster)
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	_, err = env.Simulators().
		New().
		Name(deviceID).
		Image(image).
		PullPolicy(pullPolicy).
		Add()
	return err
}

// getAddSimulatorCommand returns a cobra command for deploying a device simulator
func getAddAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app image-name [name]",
		Short: "Add an app to the test cluster",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runAddAppCommand,
	}

	cmd.Flags().StringP("image", "i", "", "the image to deploy")
	_ = cmd.MarkFlagRequired("image")
	cmd.Flags().IntP("nodes", "n", 1, "the number of nodes to deploy")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringP("cluster", "c", "", "the cluster to which to add the app")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	_ = cmd.MarkFlagRequired("cluster")
	return cmd
}

func runAddAppCommand(cmd *cobra.Command, args []string) error {
	runCommand(cmd)

	var appID string
	if len(args) > 0 {
		appID = args[0]
	}
	if appID == "" {
		appID = random.NewPetName(2)
	}

	cluster, _ := cmd.Flags().GetString("cluster")
	image, _ := cmd.Flags().GetString("image")
	nodes, _ := cmd.Flags().GetInt("nodes")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	kubeAPI, err := kube.GetAPI(cluster)
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	_, err = env.Apps().
		New().
		Name(appID).
		Nodes(nodes).
		Image(image).
		PullPolicy(pullPolicy).
		Add()
	return err
}
