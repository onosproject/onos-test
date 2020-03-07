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
	testcluster "github.com/onosproject/onos-test/pkg/cluster"
	"github.com/onosproject/onos-test/pkg/kube"
	onitcluster "github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"github.com/spf13/cobra"
)

var (
	createExample = `
		# Setup a cluster with a given name that contains one instance of each subsystem (e.g. onos-config, onos-topo)
		onit create cluster -c my-cluster

		# Setup a cluster with default name (onos) and enale onos-cli
		onit create cluster --set onos-cli.enabled=true

		# Setup a cluster that contains two instances of onos-config subsystem and two instances of onos-topo subsystem
		onit create cluster  --set onos-topo.replicas=2 --set onos-config.replicas=2

		# Setup a cluster that has two database partitions
		onit create cluster --set database.partitions=2 

		# Setup a cluster to deploy topo and config subsystems using the images with custom tags
		onit create cluster --set onos-topo.image=onosproject/onos-topo:mytag  --set onos-config.image=onosproject/onos-config:latest`
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
		RunE:  runInCluster(runCreateClusterCommand),
	}
	cmd.Flags().StringToString("set", map[string]string{}, "set a cluster argument")
	cmd.Flags().StringToString("loggers", map[string]string{}, "set required logger config files")
	return cmd
}

func runCreateClusterCommand(cmd *cobra.Command, _ []string) error {
	setupCommand(cmd)

	// Get the k8s API
	api, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	// Create the cluster
	c, err := testcluster.NewCluster(api.Namespace())
	if err != nil {
		return err
	}
	if err := c.Create(); err != nil {
		return err
	}

	args, _ := cmd.Flags().GetStringToString("set")

	onitcluster.SetArgs(args)
	setup := setup.New(api)
	setup.Atomix()
	setup.Database().Raft()
	setup.Topo().SetReplicas(1)
	setup.Config().SetReplicas(1)
	setup.RIC().SetReplicas(1)
	return setup.Setup()
}
