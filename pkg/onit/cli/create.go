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
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/onit/setup"

	corev1 "k8s.io/api/core/v1"

	"github.com/spf13/cobra"
)

var (
	createExample = `
		# Setup a cluster with a given name that contains one instance of each subsystem (e.g. onos-config, onos-topo)
		onit create cluster onit-cluster-1 

		# Setup a cluster that contains two instances of onos-config subsystem and two instances of onos-topo subsystem
		onit-create-cluster onit-cluster-2 --topo-nodes 2 --config-nodes 2

		# Setup a cluster that has two 3-node raft partitions
		onit create cluster --partitions 2 --partition-size 3

		# Setup a cluster that fetches docker images from a private docker registry
		onit create cluster --docker-registry <host>:<port>
	
		# Setup a cluster to deploy topo and config subsystems using the images with custom tags 
        onit create cluster --image-tags topo=test-topo-tag,config=test-config-tag`
)

// getCreateCommand returns a cobra "setup" command for setting up resources
func getCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create {cluster} [args]",
		Short:   "Setup a test resource on Kubernetes",
		Example: createExample,
	}
	cmd.AddCommand(getCreateClusterCommand())
	return cmd
}

// getCreateClusterCommand returns a cobra command for deploying a test cluster
func getCreateClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster [args]",
		Short: "Setup a test cluster on Kubernetes",
		Args:  cobra.NoArgs,
		RunE:  runCreateClusterCommand,
	}

	images := make(map[string]string)
	images[atomixService] = defaultAtomixImage
	images[raftService] = defaultRaftImage
	images[nopaxosSequencer] = defaultNOPaxosSequencerImage
	images[nopaxosReplica] = defaultNOPaxosReplicaImage
	images[cliService] = defaultCLIImage
	images[configService] = defaultConfigImage
	images[topoService] = defaultTopoImage

	replicas := make(map[string]int)
	replicas[atomixService] = 1
	replicas[partitionService] = 1
	replicas[cliService] = 1
	replicas[configService] = 1
	replicas[topoService] = 1

	cmd.Flags().StringToStringP("image", "i", images, "override the image to deploy for a subsystem")
	cmd.Flags().StringP("db", "d", "raft", "the database protocol to deploy (e.g. raft or nopaxos)")
	cmd.Flags().StringToIntP("replicas", "r", replicas, "set the number of replicas to deploy for a subsystem")
	cmd.Flags().IntP("partitions", "p", 1, "the number of partitions to deploy")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")

	return cmd
}

func runCreateClusterCommand(cmd *cobra.Command, _ []string) error {
	runCommand(cmd)

	images, _ := cmd.Flags().GetStringToString("image")
	replicas, _ := cmd.Flags().GetStringToInt("replicas")
	protocol, _ := cmd.Flags().GetString("db")
	partitions, _ := cmd.Flags().GetInt("partitions")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}
	cluster := cluster.New(kubeAPI)
	if err := cluster.Create(); err != nil {
		return err
	}

	setup := setup.New(kubeAPI)
	atomixImage, ok := images[atomixService]
	if !ok {
		atomixImage = defaultAtomixImage
	}
	setup.Atomix().
		SetImage(atomixImage).
		SetPullPolicy(pullPolicy)

	partitionReplicas, ok := replicas[partitionService]
	if !ok {
		partitionReplicas = 1
	}

	if protocol == raftService {
		raftImage, ok := images[raftService]
		if !ok {
			raftImage = defaultRaftImage
		}
		setup.Partitions().
			Raft().
			SetPartitions(partitions).
			SetReplicasPerPartition(partitionReplicas).
			SetImage(raftImage).
			SetPullPolicy(pullPolicy)
	} else if protocol == nopaxosService {
		replicaImage, ok := images[nopaxosReplica]
		if !ok {
			replicaImage = defaultNOPaxosReplicaImage
		}
		sequencerImage, ok := images[nopaxosSequencer]
		if !ok {
			sequencerImage = defaultNOPaxosSequencerImage
		}
		setup.Partitions().
			NOPaxos().
			SetPartitions(partitions).
			SetReplicasPerPartition(partitionReplicas).
			SetReplicaImage(replicaImage).
			SetSequencerImage(sequencerImage).
			SetPullPolicy(pullPolicy)
	} else {
		return fmt.Errorf("unknown database protocol %s", protocol)
	}

	cliImage, ok := images[cliService]
	if !ok {
		cliImage = defaultCLIImage
	}
	if replicas[cliService] > 0 {
		setup.CLI().
			SetEnabled().
			SetImage(cliImage).
			SetPullPolicy(pullPolicy)
	}

	configReplicas, ok := replicas[configService]
	if !ok {
		configReplicas = 1
	}
	configImage, ok := images[configService]
	if !ok {
		configImage = defaultConfigImage
	}
	if configReplicas > 0 {
		setup.Config().
			SetReplicas(configReplicas).
			SetImage(configImage).
			SetPullPolicy(pullPolicy)
	}

	topoReplicas, ok := replicas[topoService]
	if !ok {
		topoReplicas = 1
	}
	topoImage, ok := images[topoService]
	if !ok {
		topoImage = defaultTopoImage
	}
	if topoReplicas > 0 {
		setup.Topo().
			SetReplicas(topoReplicas).
			SetImage(topoImage).
			SetPullPolicy(pullPolicy)
	}
	return setup.Setup()
}
