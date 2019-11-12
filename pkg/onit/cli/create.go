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
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/util/random"

	"github.com/onosproject/onos-test/pkg/onit/setup"

	corev1 "k8s.io/api/core/v1"

	"github.com/spf13/cobra"
)

var (
	createExample = `
		# Create a cluster with a given name that contains one instance of each subsystem (e.g. onos-config, onos-topo)
		onit create cluster onit-cluster-1 

		# Create a cluster that contains two instances of onos-config subsystem and two instances of onos-topo subsystem
		onit-create-cluster onit-cluster-2 --topo-nodes 2 --config-nodes 2

		# Create a cluster that has two 3-node raft partitions
		onit create cluster --partitions 2 --partition-size 3

		# Create a cluster that fetches docker images from a private docker registry
		onit create cluster --docker-registry <host>:<port>
	
		# Create a cluster to deploy topo and config subsystems using the images with custom tags 
        onit create cluster --image-tags topo=test-topo-tag,config=test-config-tag`
)

// getCreateCommand returns a cobra "setup" command for setting up resources
func getCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create {cluster} [args]",
		Short:   "Create a test resource on Kubernetes",
		Example: createExample,
	}
	cmd.AddCommand(getCreateClusterCommand())
	return cmd
}

// getCreateClusterCommand returns a cobra command for deploying a test cluster
func getCreateClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster [id]",
		Short: "Setup a test cluster on Kubernetes",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runCreateClusterCommand,
	}

	images := make(map[string]string)
	images[atomixService] = defaultAtomixImage
	images[raftService] = defaultRaftImage
	images[cliService] = defaultCLIImage
	images[configService] = defaultConfigImage
	images[topoService] = defaultTopoImage

	nodes := make(map[string]int)
	nodes[atomixService] = 1
	nodes[raftService] = 1
	nodes[cliService] = 1
	nodes[configService] = 1
	nodes[topoService] = 1

	cmd.Flags().StringToStringP("image", "i", images, "override the image to deploy for a subsystem")
	cmd.Flags().StringToIntP("nodes", "n", nodes, "set the number of nodes to deploy for a subsystem")
	cmd.Flags().IntP("partitions", "p", 1, "the number of Raft partitions to deploy")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")

	return cmd
}

func runCreateClusterCommand(cmd *cobra.Command, args []string) error {
	runCommand(cmd)

	var clusterID string
	if len(args) > 0 {
		clusterID = args[0]
	}
	if clusterID == "" {
		clusterID = random.NewPetName(2)
	}

	images, _ := cmd.Flags().GetStringToString("image")
	nodes, _ := cmd.Flags().GetStringToInt("nodes")
	partitions, _ := cmd.Flags().GetInt("partitions")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	kubeAPI := kube.GetAPI(clusterID)
	cluster := cluster.New(kubeAPI)
	if err := cluster.Create(); err != nil {
		return err
	}

	setup := setup.New(kubeAPI)
	setup.Atomix().
		Image(images[atomixService]).
		PullPolicy(pullPolicy)
	setup.Database().
		Partitions(partitions).
		Nodes(nodes[raftService]).
		Image(images[raftService]).
		PullPolicy(pullPolicy)
	if nodes[cliService] > 0 {
		setup.CLI().
			Nodes(nodes[cliService]).
			Image(images[cliService]).
			PullPolicy(pullPolicy)
	}
	if nodes[configService] > 0 {
		setup.Config().
			Nodes(nodes[configService]).
			Image(images[configService]).
			PullPolicy(pullPolicy)
	}
	if nodes[topoService] > 0 {
		setup.Topo().
			Nodes(nodes[topoService]).
			Image(images[topoService]).
			PullPolicy(pullPolicy)
	}
	return setup.Setup()
}
