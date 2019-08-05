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
	"github.com/onosproject/onos-test/pkg/onit"
	"github.com/spf13/cobra"
)

var (
	deleteExample = `
            # To delete a cluster with a given name:
			onit delete cluster <name of cluster>

			# To delete the default cluster:
			onit delete cluster`
)

// getDeleteCommand returns a cobra "teardown" command for tearing down Kubernetes test resources
func getDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete Kubernetes test resources",
		Example: deleteExample,
	}
	cmd.AddCommand(getDeleteClusterCommand())
	return cmd
}

// getDeleteClusterCommand returns a cobra "teardown" command for tearing down a test cluster
func getDeleteClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster [id]",
		Short: "Delete a test cluster on Kubernetes",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Create the onit controller
			controller, err := onit.NewController()
			if err != nil {
				exitError(err)
			}

			// Get the cluster ID
			var clusterID string
			if len(args) > 0 {
				clusterID = args[0]
			} else {
				clusterID = getDefaultCluster()
			}

			// Delete the cluster
			status := controller.DeleteCluster(clusterID)
			err = setDefaultCluster("")
			if err != nil {
				exitError(err)
			}
			if status.Failed() {
				exitStatus(status)
			}
		},
	}
	return cmd
}
