// Copyright 2020-present Open Networking Foundation.
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
	"github.com/onosproject/onos-test/pkg/simulation"
	"time"

	"github.com/onosproject/onos-test/pkg/cluster"
	onitcluster "github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/util/random"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

func getSimulateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "simulate",
		Aliases: []string{"sim", "simulation"},
		Short:   "Run simulations on Kubernetes",
		RunE:    runSimulateCommand,
	}
	cmd.Flags().StringP("image", "i", "", "the simulation image to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringToString("set", map[string]string{}, "cluster argument overrides")
	cmd.Flags().StringP("simulation", "s", "", "the simulation to run")
	cmd.Flags().IntP("simulators", "w", 1, "the number of simulator workers to run")
	cmd.Flags().DurationP("rate", "r", 1*time.Second, "the rate at which to simulate operations")
	cmd.Flags().Float64P("jitter", "j", 1, "the jitter to apply to the rate")
	cmd.Flags().DurationP("duration", "d", 10*time.Minute, "the duration for which to run the simulation")
	cmd.Flags().StringToStringP("args", "a", map[string]string{}, "a mapping of named simulation arguments")

	_ = cmd.MarkFlagRequired("image")
	return cmd
}

func runSimulateCommand(cmd *cobra.Command, _ []string) error {
	setupCommand(cmd)

	image, _ := cmd.Flags().GetString("image")
	sim, _ := cmd.Flags().GetString("simulation")
	workers, _ := cmd.Flags().GetInt("simulators")
	rate, _ := cmd.Flags().GetDuration("rate")
	jitter, _ := cmd.Flags().GetFloat64("jitter")
	duration, _ := cmd.Flags().GetDuration("duration")
	sets, _ := cmd.Flags().GetStringToString("set")
	args, _ := cmd.Flags().GetStringToString("args")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	config := &simulation.Config{
		ID:              random.NewPetName(2),
		Image:           image,
		ImagePullPolicy: pullPolicy,
		Simulation:      sim,
		Simulators:      workers,
		Rate:            rate,
		Jitter:          jitter,
		Duration:        duration,
		Args:            args,
		Env:             onitcluster.GetArgsAsEnv(sets),
	}

	job := &cluster.Job{
		ID:              config.ID,
		Image:           image,
		ImagePullPolicy: pullPolicy,
		Env:             config.ToEnv(),
		Timeout:         timeout,
		Type:            "simulation",
	}

	// Create a job runner and run the benchmark job
	runner, err := cluster.NewRunner()
	if err != nil {
		return err
	}
	return runner.Run(job)
}
