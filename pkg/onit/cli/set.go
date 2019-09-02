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

	"github.com/onosproject/onos-test/pkg/onit/k8s"
	corev1 "k8s.io/api/core/v1"

	"github.com/spf13/cobra"
)

var (
	setExample = ` 
		# Change the currently configured cluster
		onit set cluster <name of a cluster>
		
		# Update existing container image(s) of deployments.
		onit set image onos-config --image onosproject/onos-config:debug`
)

// getSetCommand returns a cobra "set" command for setting configurations
func getSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set {cluster,image} [args]",
		Short:   "Set test configurations",
		Example: setExample,
	}
	cmd.AddCommand(getSetClusterCommand())
	cmd.AddCommand(getSetImageCommand())
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
			controller, err := k8s.NewController()
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

// getSetImageCommand returns a cobra command for update existing container image(s) of resources.
func getSetImageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image <deployment_name>",
		Args:  cobra.ExactArgs(1),
		Short: "Set image <deployment_name>",
		Run: func(cmd *cobra.Command, args []string) {
			nodeID := args[0]
			image, _ := cmd.Flags().GetString("image")
			imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
			pullPolicy := corev1.PullPolicy(imagePullPolicy)

			if pullPolicy != corev1.PullAlways && pullPolicy != corev1.PullIfNotPresent && pullPolicy != corev1.PullNever {
				exitError(fmt.Errorf("invalid pull policy; must of one of %s, %s or %s", corev1.PullAlways, corev1.PullIfNotPresent, corev1.PullNever))
			}

			clusterID, err := cmd.Flags().GetString("cluster")
			if err != nil {
				exitError(err)
			}

			// Get the onit controller
			controller, err := k8s.NewController()
			if err != nil {
				exitError(err)
			}

			// Get the cluster controller
			cluster, err := controller.GetCluster(clusterID)
			if err != nil {
				exitError(err)
			}

			status := cluster.SetImage(nodeID, image, pullPolicy)
			if status.Failed() {
				exitStatus(status)
			}
		},
	}
	cmd.Flags().StringP("cluster", "c", getDefaultCluster(), "the cluster to query")
	cmd.Flags().StringP("image", "i", "", "the image name")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	_ = cmd.MarkFlagRequired("image")
	return cmd
}
