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
	"fmt"
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"os"
	"strconv"
	"strings"
	"time"
)

type simulationContext string

const (
	simulationContextEnv = "SIMULATION_CONTEXT"

	simulationJobEnv             = "SIMULATION_JOB"
	simulationImageEnv           = "SIMULATION_IMAGE"
	simulationImagePullPolicyEnv = "SIMULATION_IMAGE_PULL_POLICY"
	simulationNameEnv            = "SIMULATION_NAME"
	simulationModelEnv           = "SIMULATION_MODEL"
	simulationSimulatorsEnv      = "SIMULATION_SIMULATORS"
	simulationRateEnv            = "SIMULATION_RATE"
	simulationJitterEnv          = "SIMULATION_JITTER"
	simulationDurationEnv        = "SIMULATION_DURATION"
	simulationParallelismEnv     = "SIMULATION_PARALLELISM"
	simulationArgsEnv            = "SIMULATION_ARGS"
	simulationWorkerEnv          = "SIMULATION_WORKER"
)

const (
	simulationContextCoordinator simulationContext = "coordinator"
	simulationContextWorker      simulationContext = "worker"
)

// GetConfigFromEnv returns the simulation configuration from the environment
func GetConfigFromEnv() *Config {
	env := make(map[string]string)
	for _, keyval := range os.Environ() {
		key := keyval[:strings.Index(keyval, "=")]
		value := keyval[strings.Index(keyval, "=")+1:]
		env[key] = value
	}
	args := make(map[string]string)
	for key, value := range cluster.GetArgs() {
		args[key] = value
	}
	for key, value := range util.SplitMap(os.Getenv(simulationArgsEnv)) {
		args[key] = value
	}
	workers, err := strconv.Atoi(os.Getenv(simulationSimulatorsEnv))
	if err != nil {
		panic(err)
	}
	var rate time.Duration
	rateEnv := os.Getenv(simulationRateEnv)
	if rateEnv != "" {
		d, err := strconv.Atoi(rateEnv)
		if err != nil {
			panic(err)
		}
		rate = time.Duration(d)
	}
	jitter, err := strconv.ParseFloat(os.Getenv(simulationJitterEnv), 64)
	if err != nil {
		panic(err)
	}
	var duration time.Duration
	durationEnv := os.Getenv(simulationDurationEnv)
	if durationEnv != "" {
		d, err := strconv.Atoi(durationEnv)
		if err != nil {
			panic(err)
		}
		duration = time.Duration(d)
	}
	parallelism, err := strconv.Atoi(os.Getenv(simulationParallelismEnv))
	if err != nil {
		panic(err)
	}
	return &Config{
		ID:              os.Getenv(simulationJobEnv),
		Image:           os.Getenv(simulationImageEnv),
		ImagePullPolicy: corev1.PullPolicy(os.Getenv(simulationImagePullPolicyEnv)),
		Simulation:      os.Getenv(simulationNameEnv),
		Model:           os.Getenv(simulationModelEnv),
		Simulators:      workers,
		Rate:            rate,
		Jitter:          jitter,
		Duration:        duration,
		Parallelism:     parallelism,
		Args:            args,
		Env:             env,
	}
}

// Config is a simulation configuration
type Config struct {
	ID              string
	Image           string
	ImagePullPolicy corev1.PullPolicy
	Simulation      string
	Model           string
	Simulators      int
	Rate            time.Duration
	Jitter          float64
	Parallelism     int
	Duration        time.Duration
	Args            map[string]string
	Env             map[string]string
}

// ToEnv returns the configuration as a mapping of environment variables
func (c *Config) ToEnv() map[string]string {
	env := c.Env
	env[simulationJobEnv] = c.ID
	env[simulationImageEnv] = c.Image
	env[simulationImagePullPolicyEnv] = string(c.ImagePullPolicy)
	env[simulationNameEnv] = c.Simulation
	env[simulationModelEnv] = c.Model
	env[simulationSimulatorsEnv] = fmt.Sprintf("%d", c.Simulators)
	env[simulationRateEnv] = fmt.Sprintf("%d", c.Rate)
	env[simulationJitterEnv] = fmt.Sprintf("%f", c.Jitter)
	env[simulationDurationEnv] = fmt.Sprintf("%d", c.Duration)
	env[simulationParallelismEnv] = fmt.Sprintf("%d", c.Parallelism)
	env[simulationArgsEnv] = util.JoinMap(c.Args)
	return env
}

// getSimulationContext returns the current simulation context
func getSimulationContext() simulationContext {
	context := os.Getenv(simulationContextEnv)
	if context != "" {
		return simulationContext(context)
	}
	return simulationContextCoordinator
}

// getSimulatorID returns the current simulation worker number
func getSimulatorID() int {
	worker := os.Getenv(simulationWorkerEnv)
	if worker == "" {
		return 0
	}
	i, err := strconv.Atoi(worker)
	if err != nil {
		panic(err)
	}
	return i
}
