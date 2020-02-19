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
	"net/http"
	"net/url"
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

var (
	simulateExample = `
		# Simulate operations on an Atomix map
		onit simulate --image atomix/kubernetes-simulations --simulation map --duration 1m

		# Configure the simulated Atomix cluster
		onit simulate --image atomix/kubernetes-simulations --simulation map --duration 1m --set raft.clusters=3 --set raft.partitions=3

		# Configure scheduled operations on an Atomix map
		onit simulate --image atomix/kubernetes-simulations --simulation map --schedule put=2s --schedule get=1s,.5 --schedule remove=5s --duration 5m

		# Verify an Atomix map simulation against a TLA+ model
		onit simulate --image atomix/kubernetes-simulations --simulation map --duration 5m --verify --model models/MapCacheTrace.tla --module models/MapHistory.tla --spec Spec --invariant StateInvariant`
)

func getSimulateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "simulate",
		Aliases: []string{"sim", "simulation"},
		Short:   "Run simulations on Kubernetes",
		Example: simulateExample,
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

// isURL returns whether the given string is a URL
func isURL(str string) bool {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

// downloadURL downloads the given URL to a []byte
func downloadURL(url string) ([]byte, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// getData gets the name and data for the given file system path or URL
func getData(pathOrURL string) (string, []byte, error) {
	var name string
	var data []byte
	if isURL(pathOrURL) {
		bytes, err := downloadURL(pathOrURL)
		if err != nil {
			return "", nil, err
		}
		u, err := url.Parse(pathOrURL)
		if err != nil {
			return "", nil, err
		}
		name = path.Base(u.Path)
		name = name[:len(name)-len(path.Ext(name))]
		data = bytes
	} else {
		file, err := os.Open(pathOrURL)
		if err != nil {
			return "", nil, err
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return "", nil, err
		}
		name = path.Base(pathOrURL)
		name = name[:len(name)-len(path.Ext(name))]
		data = bytes
	}
	return name, data, nil
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

		modelData = make(map[string]string)
		if modelPath != "" {
			name, data, err := getData(modelPath)
			if err != nil {
				return err
			}
			modelName = name
			modelData[fmt.Sprintf("%s.tla", name)] = string(data)

			var configBytes []byte
			if configPath != "" {
				_, data, err := getData(configPath)
				if err != nil {
					return err
				}
				configBytes = data
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
			modelData[fmt.Sprintf("%s.cfg", modelName)] = string(configBytes)

			for _, modulePath := range modulePaths {
				name, data, err := getData(modulePath)
				if err != nil {
					return err
				}
				modelData[fmt.Sprintf("%s.tla", name)] = string(data)
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
