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
	"bytes"
	"fmt"
	"github.com/onosproject/onos-test/pkg/simulation"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
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
	cmd.Flags().DurationP("duration", "d", 10*time.Minute, "the duration for which to run the simulation")
	cmd.Flags().StringToStringP("args", "a", map[string]string{}, "a mapping of named simulation arguments")
	cmd.Flags().StringToStringP("schedule", "r", map[string]string{}, "a mapping of operations to schedule")
	cmd.Flags().Bool("verify", false, "whether to verify the simulation against a formal model")
	cmd.Flags().StringP("model", "m", "", "a model with which to verify the simulation")
	cmd.Flags().StringArray("module", []string{}, "modules to add to the model")
	cmd.Flags().String("config", "", "the model configuration")
	cmd.Flags().String("spec", "", "the model specification")
	cmd.Flags().String("init", "", "an init predicate")
	cmd.Flags().String("next", "", "a next state predicate")
	cmd.Flags().StringArray("invariant", []string{}, "model invariant")
	cmd.Flags().StringArray("constant", []string{}, "model constants")
	cmd.Flags().StringArray("constraint", []string{}, "model constraints")
	cmd.Flags().StringArray("property", []string{}, "model properties")

	_ = cmd.MarkFlagRequired("image")
	return cmd
}

func runSimulateCommand(cmd *cobra.Command, _ []string) error {
	setupCommand(cmd)

	image, _ := cmd.Flags().GetString("image")
	sim, _ := cmd.Flags().GetString("simulation")
	workers, _ := cmd.Flags().GetInt("simulators")
	duration, _ := cmd.Flags().GetDuration("duration")
	sets, _ := cmd.Flags().GetStringToString("set")
	args, _ := cmd.Flags().GetStringToString("args")
	operations, _ := cmd.Flags().GetStringToString("schedule")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	verify, _ := cmd.Flags().GetBool("verify")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	var modelName string
	var modelData map[string]string
	if verify {
		modelPath, _ := cmd.Flags().GetString("model")
		modulePaths, _ := cmd.Flags().GetStringArray("module")
		configPath, _ := cmd.Flags().GetString("config")

		if modelPath != "" {
			modelName = path.Base(modelPath)
			modelName = modelName[:len(modelName)-len(path.Ext(modelName))]

			var configBytes []byte
			if configPath != "" {
				file, err := os.Open(configPath)
				if err != nil {
					return err
				}
				bytes, err := ioutil.ReadAll(file)
				if err != nil {
					return err
				}
				configBytes = bytes
			} else {
				spec, _ := cmd.Flags().GetString("spec")
				init, _ := cmd.Flags().GetString("init")
				next, _ := cmd.Flags().GetString("next")
				invariants, _ := cmd.Flags().GetStringArray("invariant")
				constants, _ := cmd.Flags().GetStringArray("constant")
				constraints, _ := cmd.Flags().GetStringArray("constraint")
				properties, _ := cmd.Flags().GetStringArray("property")

				buf := &bytes.Buffer{}
				if spec != "" {
					fmt.Fprintln(buf, "SPECIFICATION", spec)
				}
				if init != "" {
					fmt.Fprintln(buf, "INIT", init)
				}
				if next != "" {
					fmt.Fprintln(buf, "NEXT", next)
				}
				if len(invariants) > 0 {
					for _, invariant := range invariants {
						fmt.Fprintln(buf, "INVARIANT", invariant)
					}
				}
				if len(constants) > 0 {
					for _, constant := range constants {
						fmt.Fprintln(buf, "CONSTANT", constant)
					}
				}
				if len(constraints) > 0 {
					for _, constraint := range constraints {
						fmt.Fprintln(buf, "CONSTRAINT", constraint)
					}
				}
				if len(properties) > 0 {
					for _, property := range properties {
						fmt.Fprintln(buf, "PROPERTY", property)
					}
				}
				configBytes = buf.Bytes()
			}

			modelData = make(map[string]string)
			file, err := os.Open(modelPath)
			if err != nil {
				return err
			}
			modelBytes, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			modelData[fmt.Sprintf("%s.tla", modelName)] = string(modelBytes)
			modelData[fmt.Sprintf("%s.cfg", modelName)] = string(configBytes)

			for _, modulePath := range modulePaths {
				file, err := os.Open(modulePath)
				if err != nil {
					return err
				}
				moduleBytes, err := ioutil.ReadAll(file)
				if err != nil {
					return err
				}
				modelData[path.Base(modulePath)] = string(moduleBytes)
			}
		}
	}

	rates := make(map[string]time.Duration)
	jitters := make(map[string]float64)
	for name, value := range operations {
		var rate string
		index := strings.Index(value, ",")
		if index == -1 {
			rate = value
		} else {
			rate = value[:index]
			jitter := value[index+1:]
			f, err := strconv.ParseFloat(jitter, 64)
			if err != nil {
				return err
			}
			jitters[name] = f
		}
		d, err := time.ParseDuration(rate)
		if err != nil {
			return err
		}
		rates[name] = d
	}

	config := &simulation.Config{
		ID:              random.NewPetName(2),
		Image:           image,
		ImagePullPolicy: pullPolicy,
		Simulation:      sim,
		Model:           modelName,
		Simulators:      workers,
		Duration:        duration,
		Rates:           rates,
		Jitter:          jitters,
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
		ModelChecker:    verify,
		ModelData:       modelData,
	}

	// Create a job runner and run the benchmark job
	runner, err := cluster.NewRunner()
	if err != nil {
		return err
	}
	return runner.Run(job)
}
