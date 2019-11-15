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
	"strings"

	"github.com/spf13/cobra"
)

var (
	execExample = `
        # Execute a command on the ONOS CLI
        onit exec -- onos topo get devices`
)

// getExecCommand returns a cobra "exec" command for executing CLI commands
func getExecCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "exec {command}",
		Short:   "Execute an ONOS CLI command",
		Example: execExample,
		RunE:    runExecCommand,
	}
}

// runExecCommand runs the "exec" command
func runExecCommand(cmd *cobra.Command, args []string) error {
	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	output, code, err := env.CLI().Execute(strings.Join(args, " "))
	if err != nil {
		return err
	}

	for _, line := range output {
		fmt.Println(line)
	}

	os.Exit(code)
	return nil
}
