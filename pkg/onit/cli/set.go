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
	setExample = ` 
		# Change the currently configured cluster
		onit set cluster <name of a cluster>`
)

// getSetCommand returns a cobra "set" command for setting configurations
func getSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set {cluster} [args]",
		Short:   "Set test configurations",
		Example: setExample,
	}
	cmd.AddCommand(getSetClusterCommand())
	return cmd
}

// getSetClusterCommand returns a cobra command for setting the cluster context
func getSetClusterCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cluster <name>",
		Args:  cobra.ExactArgs(1),
		Short: "Set cluster context",
		Run: func(cmd *cobra.Command, args []string) {
			clusterID := args[0]

			// Get the onit controller
			controller, err := onit.NewController()
			if err != nil {
				exitError(err)
			}

			// Get the cluster controller
			_, err = controller.GetCluster(clusterID)
			if err != nil {
				exitError(err)
			}

			// If we've made it this far, update the default cluster
			if err := setDefaultCluster(clusterID); err != nil {
				exitError(err)
			} else {
				fmt.Println(clusterID)
			}
		},
	}
}
