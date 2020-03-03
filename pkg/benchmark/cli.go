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

package benchmark

import (
	"github.com/onosproject/onos-test/pkg/cluster"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"math/rand"
	"time"
)

// GetCommand returns the test command
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kube-bench",
		Short: "Start and manage Kubernetes benchmarks",
		RunE:  runBenchCommand,
	}
	cmd.Flags().StringP("image", "i", "", "the test image to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the image pull policy to use")
	cmd.Flags().StringP("suite", "s", "", "the benchmark suite to run")
	cmd.Flags().StringP("benchmark", "b", "", "the name of the benchmark to run")
	cmd.Flags().IntP("workers", "w", 1, "the number of workers to run")
	cmd.Flags().IntP("parallel", "p", 1, "the number of concurrent goroutines per client")
	cmd.Flags().IntP("requests", "n", 1, "the number of requests to run")
	cmd.Flags().StringToStringP("args", "a", map[string]string{}, "a mapping of named benchmark arguments")
	cmd.Flags().Duration("timeout", 0, "the test timeout")
	return cmd
}

// runBenchCommand runs the bench command
func runBenchCommand(cmd *cobra.Command, _ []string) error {
	image, _ := cmd.Flags().GetString("image")
	pullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	suite, _ := cmd.Flags().GetString("suite")
	benchmark, _ := cmd.Flags().GetString("benchmark")
	workers, _ := cmd.Flags().GetInt("workers")
	parallelism, _ := cmd.Flags().GetInt("parallel")
	requests, _ := cmd.Flags().GetInt("requests")
	args, _ := cmd.Flags().GetStringToString("args")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	maxLatency, _ := cmd.Flags().GetDuration("max-latency")

	config := &Config{
		ID:              random.NewPetName(2),
		Image:           image,
		ImagePullPolicy: corev1.PullPolicy(pullPolicy),
		Suite:           suite,
		Benchmark:       benchmark,
		Workers:         workers,
		Parallelism:     parallelism,
		Requests:        requests,
		Args:            args,
		Timeout:         timeout,
		MaxLatency:      &maxLatency,
	}

	job := &cluster.Job{
		ID:              config.ID,
		Image:           image,
		ImagePullPolicy: corev1.PullPolicy(pullPolicy),
		Env:             config.ToEnv(),
		Timeout:         timeout,
	}
	runner, err := cluster.NewRunner()
	if err != nil {
		return err
	}
	return runner.Run(job)
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
