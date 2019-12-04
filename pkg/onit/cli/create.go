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
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"github.com/spf13/cobra"
)

var (
	createExample = `
		# Setup a cluster with a default name (i.e. onos) that contains one instance of each subsystem (e.g. onos-config, onos-topo, atomix controller, database)
		onit create cluster  

		# Setup a cluster with a given name that contains two instances of onos-config subsystem and two instances of onos-topo subsystem
		onit create cluster -c onit-cluster-1 --set onos-topo.replicas=2 --set onos-config.replicas=2
		
		# Setup a cluster with 3 
`
)

// getCreateCommand returns a cobra "setup" command for setting up resources
func getCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create {cluster} [options]",
		Short:   "Setup a test resource on Kubernetes",
		Example: createExample,
	}
	cmd.AddCommand(getCreateClusterCommand())
	return cmd
}

// getCreateClusterCommand returns a cobra command for deploying a test cluster
func getCreateClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster [options]",
		Short: "Setup a test cluster on Kubernetes",
		Args:  cobra.NoArgs,
		RunE:  runCreateClusterCommand,
	}
	cmd.Flags().StringToString("set", map[string]string{}, "set a cluster argument")
	return cmd
}

func runCreateClusterCommand(cmd *cobra.Command, _ []string) error {
	runCommand(cmd)
	args, _ := cmd.Flags().GetStringToString("set")
	cluster.SetArgs(args)
	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}
	cluster := cluster.New(kubeAPI)
	if err := cluster.Create(); err != nil {
		return err
	}
	setup := setup.New(kubeAPI)
	setup.Atomix()
	setup.Partitions().Raft()
	setup.Topo().SetReplicas(1)
	setup.Config().SetReplicas(1)
	return setup.Setup()
}
