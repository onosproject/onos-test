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
	"github.com/onosproject/onos-test/pkg/test"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"os"
	"strings"
	"time"
)

func getTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "test",
		Aliases: []string{"tests"},
		Short:   "Run tests on Kubernetes",
		RunE:    runTestCommand,
	}
	cmd.Flags().StringP("image", "i", "", "the test image to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringToString("set", map[string]string{}, "cluster argument overrides")
	cmd.Flags().StringP("suite", "s", "", "the test suite to run")
	cmd.Flags().StringP("test", "t", "", "the name of the test method to run")
	cmd.Flags().Duration("timeout", 10*time.Minute, "test timeout")
	cmd.Flags().Bool("no-teardown", false, "do not tear down clusters following tests")
	return cmd
}

func runTestCommand(cmd *cobra.Command, _ []string) error {
	runCommand(cmd)

	clusterID, _ := cmd.Flags().GetString("cluster")
	image, _ := cmd.Flags().GetString("image")
	sets, _ := cmd.Flags().GetStringToString("set")
	suite, _ := cmd.Flags().GetString("suite")
	testName, _ := cmd.Flags().GetString("test")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	noTeardown, _ := cmd.Flags().GetBool("no-teardown")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	overrides := []string{}
	for key, value := range sets {
		overrides = append(overrides, fmt.Sprintf("%s=%s", key, value))
	}

	config := &test.TestConfig{
		JobConfig: &test.JobConfig{
			JobID:      random.NewPetName(2),
			Type:       test.TestTypeTest,
			Image:      image,
			Env: map[string]string{
				"ARGS": strings.Join(overrides, ","),
			},
			Timeout:    timeout,
			PullPolicy: corev1.PullPolicy(pullPolicy),
			Teardown:   !noTeardown,
		},
		Suite: suite,
		Test:  testName,
	}

	// If the cluster ID was not specified, create a new cluster to run the test
	// Otherwise, deploy the test in the existing cluster
	if clusterID == "" {
		runner, err := test.NewTestRunner(config)
		if err != nil {
			return err
		}
		return runner.Run()
	}

	cluster, err := test.NewTestCluster(clusterID)
	if err != nil {
		return err
	}
	if err := cluster.StartTest(config); err != nil {
		return err
	}
	if err := cluster.AwaitTestComplete(config); err != nil {
		return err
	}
	_, code, err := cluster.GetTestResult(config)
	if err != nil {
		return err
	}
	os.Exit(code)
	return nil
}
