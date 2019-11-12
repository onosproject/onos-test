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
	"github.com/onosproject/onos-test/pkg/onit/env"
	"os"

	"github.com/spf13/cobra"
)

var (
	execExample = `
        # Execute a command on the ONOS CLI
        onit exec -- onos topo get devices`
)

// getExecCommand returns a cobra "exec" command for executing CLI commands
func getExecCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exec {command}",
		Short:   "Execute an ONOS CLI command",
		Example: execExample,
		RunE:    runExecCommand,
	}
	cmd.Flags().StringP("cluster", "c", "", "the cluster to which to add the simulator")
	cmd.Flags().Lookup("cluster").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__onit_get_clusters"},
	}
	_ = cmd.MarkFlagRequired("cluster")
	return cmd
}

// runExecCommand runs the "exec" command
func runExecCommand(cmd *cobra.Command, args []string) error {
	cluster, _ := cmd.Flags().GetString("cluster")
	kubeAPI := kube.GetAPI(cluster)
	env := env.New(kubeAPI)
	output, code, err := env.CLI().Execute(args...)
	if err != nil {
		return err
	}

	for _, line := range output {
		fmt.Println(line)
	}

	os.Exit(code)
	return nil
}
