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
	"github.com/onosproject/onos-test/pkg/cluster"
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
	simulationDurationEnv        = "SIMULATION_DURATION"
	simulationRatesEnv           = "SIMULATION_RATES"
	simulationJittersEnv         = "SIMULATION_JITTERS"
	simulationArgsEnv            = "SIMULATION_ARGS"
	simulationWorkerEnv          = "SIMULATION_WORKER"
)

const (
	simulationContextCoordinator simulationContext = "coordinator"
	simulationContextWorker      simulationContext = "worker"
)

// getAddress returns the service address
func getAddress() string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:5000", os.Getenv("SERVICE_NAME"), os.Getenv("SERVICE_NAMESPACE"))
}

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
	rates := make(map[string]time.Duration)
	for key, value := range util.SplitMap(os.Getenv(simulationRatesEnv)) {
		d, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		rates[key] = time.Duration(d)
	}
	jitter := make(map[string]float64)
	for key, value := range util.SplitMap(os.Getenv(simulationJittersEnv)) {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(err)
		}
		jitter[key] = f
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
	return &Config{
		ID:              os.Getenv(simulationJobEnv),
		Image:           os.Getenv(simulationImageEnv),
		ImagePullPolicy: corev1.PullPolicy(os.Getenv(simulationImagePullPolicyEnv)),
		Simulation:      os.Getenv(simulationNameEnv),
		Model:           os.Getenv(simulationModelEnv),
		Simulators:      workers,
		Duration:        duration,
		Rates:           rates,
		Jitter:          jitter,
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
	Rates           map[string]time.Duration
	Jitter          map[string]float64
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
	env[simulationDurationEnv] = fmt.Sprintf("%d", c.Duration)
	rates := make(map[string]string)
	for key, rate := range c.Rates {
		rates[key] = fmt.Sprintf("%d", rate)
	}
	env[simulationRatesEnv] = util.JoinMap(rates)
	jitters := make(map[string]string)
	for key, jitter := range c.Jitter {
		jitters[key] = strconv.FormatFloat(jitter, 'E', -1, 64)
	}
	env[simulationJittersEnv] = util.JoinMap(jitters)
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
