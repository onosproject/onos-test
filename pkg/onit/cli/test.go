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
	"github.com/onosproject/onos-test/pkg/cluster"
	onitcluster "github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/test"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"time"
)

func getTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "test",
		Aliases: []string{"tests"},
		Short:   "Run tests on Kubernetes",
		RunE:    runTestCommand,
	}
	defaultSlice := make([]string, 0)
	cmd.Flags().StringP("image", "i", "", "the test image to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringToString("set", map[string]string{}, "cluster argument overrides")
	cmd.Flags().StringSliceP("suite", "s", defaultSlice, "the name of test suite to run")
	cmd.Flags().StringSliceP("test", "t", defaultSlice, "the name of the test method to run")
	cmd.Flags().Duration("timeout", 10*time.Minute, "test timeout")
	cmd.Flags().Int("iterations", 1, "number of iterations")
	cmd.Flags().Bool("until-failure", false, "run until an error is detected")
	cmd.Flags().Bool("no-teardown", false, "do not tear down clusters following tests")

	_ = cmd.MarkFlagRequired("image")
	return cmd
}

func runTestCommand(cmd *cobra.Command, _ []string) error {
	setupCommand(cmd)

	image, _ := cmd.Flags().GetString("image")
	sets, _ := cmd.Flags().GetStringToString("set")
	suites, _ := cmd.Flags().GetStringSlice("suite")
	testNames, _ := cmd.Flags().GetStringSlice("test")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	pullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	iterations, _ := cmd.Flags().GetInt("iterations")
	untilFailure, _ := cmd.Flags().GetBool("until-failure")

	if untilFailure {
		iterations = -1
	}

	config := &test.Config{
		ID:              random.NewPetName(2),
		Image:           image,
		ImagePullPolicy: corev1.PullPolicy(pullPolicy),
		Suites:          suites,
		Tests:           testNames,
		Env:             onitcluster.GetArgsAsEnv(sets),
		Timeout:         timeout,
		Iterations:      iterations,
	}

	job := &cluster.Job{
		ID:              config.ID,
		Image:           image,
		ImagePullPolicy: corev1.PullPolicy(pullPolicy),
		Env:             config.ToEnv(),
		Timeout:         timeout,
	}

	// Create a job runner and run the test job
	runner, err := cluster.NewRunner()
	if err != nil {
		return err
	}
	return runner.Run(job)
}
