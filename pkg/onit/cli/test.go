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
	"github.com/onosproject/onos-test/pkg/kubetest"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"os"
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
	suite, _ := cmd.Flags().GetString("suite")
	test, _ := cmd.Flags().GetString("test")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	noTeardown, _ := cmd.Flags().GetBool("no-teardown")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	config := &kubetest.TestConfig{
		TestID:     random.NewPetName(2),
		Type:       kubetest.TestTypeTest,
		Image:      image,
		Suite:      suite,
		Test:       test,
		Timeout:    timeout,
		PullPolicy: pullPolicy,
		Teardown:   !noTeardown,
	}

	// If the cluster ID was not specified, create a new cluster to run the test
	// Otherwise, deploy the test in the existing cluster
	if clusterID == "" {
		runner, err := kubetest.NewTestRunner(config)
		if err != nil {
			return err
		}
		return runner.Run()
	}

	cluster, err := kubetest.NewTestCluster(clusterID)
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
