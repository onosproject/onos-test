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
	"github.com/onosproject/onos-test/pkg/onit/k8s"
	"github.com/spf13/cobra"
)

var (
	cliExample = `
	# open an ssh connection to onos-cli node in the cluster to execute remote commands
	onit onos-cli [command args]...`
	sshExample = `
	# open an ssh connection to the specified node in the cluster to execute remote commands
	onit ssh <name of a node> [command args]...`
)

func getOnosCliCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "onos-cli",
		Short:   "Open onos-cli shell for executing commands",
		Example: cliExample,
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

			onosCliNodes, err := cluster.GetOnosCliNodes()
			if err != nil {
				exitError(err)
			}

			err = cluster.OpenShell(onosCliNodes[0].ID, args...)
			if err != nil {
				exitError(err)
			}

		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster for which to run onos-cli")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd

}

// getSshCommand returns a cobra "ssh" command to open a ssh session to a node for executing remote commands
func getSSHCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh <resource>",
		Short:   "Open a ssh session to a node for executing remote commands",
		Example: sshExample,
		Args:    cobra.ExactArgs(1),
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

			err = cluster.OpenShell(args[0], args[1:]...)
			if err != nil {
				exitError(err)
			}

		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster for which to ssh into nodes")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	return cmd
}
