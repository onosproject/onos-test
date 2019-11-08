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
	"github.com/onosproject/onos-test/pkg/new/kubetest"
	"github.com/onosproject/onos-test/pkg/new/onit/env"
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"github.com/spf13/cobra"
)

var (
	getExample = ` 
		# Get a list of clusters
		onit get clusters

		# Get the list of installed apps
		onit get apps

		# Get a list of simulators
		onit get simulators

		# Get a list of networks
		onit get networks`
)

// getGetCommand returns a cobra "get" command to read test configurations
func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get {cluster,clusters,networks,simulators}",
		Short:   "Get test configurations",
		Example: getExample,
	}
	cmd.AddCommand(getGetClustersCommand())
	cmd.AddCommand(getGetSimulatorsCommand())
	cmd.AddCommand(getGetNetworksCommand())
	cmd.AddCommand(getGetAppsCommand())
	return cmd
}

// getGetClustersCommand returns a cobra command to get a list of available test clusters
func getGetClustersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clusters",
		Short: "Get a list of all deployed clusters",
		RunE:  runGetClustersCommand,
	}
}

func runGetClustersCommand(cmd *cobra.Command, args []string) error {
	clusters, err := kubetest.GetTestClusters()
	if err != nil {
		return err
	}

	for _, cluster := range clusters {
		fmt.Println(cluster)
	}
	return nil
}

// getGetNetworksCommand returns a cobra command to get the list of networks deployed in the current cluster context
func getGetNetworksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "networks",
		Short: "Get the currently configured cluster's networks",
		RunE:  runGetNetworksCommand,
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to query")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	_ = cmd.MarkFlagRequired("cluster")
	return cmd
}

func runGetNetworksCommand(cmd *cobra.Command, _ []string) error {
	cluster, _ := cmd.Flags().GetString("cluster")
	kubeAPI := kube.GetAPI(cluster)
	env := env.New(kubeAPI)
	for _, network := range env.Networks().List() {
		fmt.Println(network.Name())
	}
	return nil
}

// getGetSimulatorsCommand returns a cobra command to get the list of simulators deployed in the current cluster context
func getGetSimulatorsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulators",
		Short: "Get the currently configured cluster's simulators",
		RunE:  runGetSimulatorsCommand,
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to query")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	_ = cmd.MarkFlagRequired("cluster")
	return cmd
}

func runGetSimulatorsCommand(cmd *cobra.Command, _ []string) error {
	cluster, _ := cmd.Flags().GetString("cluster")
	kubeAPI := kube.GetAPI(cluster)
	env := env.New(kubeAPI)
	for _, simulator := range env.Simulators().List() {
		fmt.Println(simulator.Name())
	}
	return nil
}

// getGetAppsCommand returns a cobra command to get the list of apps deployed in the current cluster context
func getGetAppsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: "Get the currently configured cluster's apps",
		RunE:  runGetAppsCommand,
	}

	cmd.Flags().StringP("cluster", "c", setup.GetDefaultCluster(), "the cluster to query")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	_ = cmd.MarkFlagRequired("cluster")
	return cmd
}

func runGetAppsCommand(cmd *cobra.Command, _ []string) error {
	cluster, _ := cmd.Flags().GetString("cluster")
	kubeAPI := kube.GetAPI(cluster)
	env := env.New(kubeAPI)
	for _, app := range env.Apps().List() {
		fmt.Println(app.Name())
	}
	return nil
}
