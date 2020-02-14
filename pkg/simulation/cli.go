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

package simulation

import (
	"github.com/onosproject/onos-test/pkg/cluster"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"math/rand"
	"time"
)

// GetCommand returns the simulate command
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kube-simulate",
		Short: "Start and manage Kubernetes simulations",
		RunE:  runSimulateCommand,
	}
	cmd.Flags().StringP("image", "i", "", "the simulation image to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the image pull policy to use")
	cmd.Flags().StringP("simulation", "s", "", "the simulation suite to run")
	cmd.Flags().IntP("simulators", "w", 1, "the number of simulator workers to run")
	cmd.Flags().DurationP("rate", "r", 1*time.Second, "the rate at which to simulate operations")
	cmd.Flags().Float64P("jitter", "j", 1, "the jitter to apply to the rate")
	cmd.Flags().DurationP("duration", "d", 10*time.Minute, "the duration for which to run the simulation")
	cmd.Flags().StringToStringP("args", "a", map[string]string{}, "a mapping of named simulation arguments")
	return cmd
}

// runSimulateCommand runs the simulate command
func runSimulateCommand(cmd *cobra.Command, _ []string) error {
	image, _ := cmd.Flags().GetString("image")
	pullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	simulation, _ := cmd.Flags().GetString("simulation")
	simulators, _ := cmd.Flags().GetInt("simulators")
	rate, _ := cmd.Flags().GetDuration("rate")
	jitter, _ := cmd.Flags().GetFloat64("jitter")
	duration, _ := cmd.Flags().GetDuration("duration")
	args, _ := cmd.Flags().GetStringToString("args")

	config := &Config{
		ID:              random.NewPetName(2),
		Image:           image,
		ImagePullPolicy: corev1.PullPolicy(pullPolicy),
		Simulation:      simulation,
		Simulators:      simulators,
		Rate:            rate,
		Jitter:          jitter,
		Duration:        duration,
		Args:            args,
	}

	job := &cluster.Job{
		ID:              config.ID,
		Image:           image,
		ImagePullPolicy: corev1.PullPolicy(pullPolicy),
		Env:             config.ToEnv(),
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
