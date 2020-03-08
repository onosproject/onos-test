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
	"os"

	"github.com/onosproject/onos-test/pkg/cluster"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

const (
	onitImage      = "onosproject/onit:latest"
	onitPullPolicy = corev1.PullIfNotPresent
)

// cliContext indicates the context in which the CLI is running
type cliContext string

const cliContextEnv = "ONIT_CONTEXT"

const (
	localContext cliContext = "local"
	k8sContext   cliContext = "kubernetes"
)

// getContext returns the current CLI context
func getContext() cliContext {
	context := os.Getenv(cliContextEnv)
	if context == "" {
		return localContext
	}
	return cliContext(context)
}

// runInCluster wraps the given command to execute it inside the cluster
func runInCluster(command func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// If running inside k8s, simply run the command
		if getContext() == k8sContext {
			return command(cmd, args)
		}

		timeout, err := cmd.Flags().GetDuration("timeout")
		if err != nil {
			timeout = 0
		}

		config, err := cmd.Flags().GetStringSlice("config")
		if err != nil {
			config = []string{}
		}

		runner, err := cluster.NewRunner()
		if err != nil {
			return err
		}

		job := &cluster.Job{
			ID:              fmt.Sprintf("onit-%s", random.NewPetName(2)),
			Image:           onitImage,
			ImagePullPolicy: onitPullPolicy,
			Args:            os.Args[1:],
			Env: map[string]string{
				cliContextEnv: string(k8sContext),
			},
			Timeout: timeout,
			Config:  config,
		}
		return runner.Run(job)
	}
}
