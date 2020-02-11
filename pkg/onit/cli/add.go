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
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/spf13/cobra"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
)

var (
	addExample = `
		# Add a simulator with a given name to "onos" cluster
		onit add simulator -n simulator-1

		# Add a simulator with a random name to a specific cluster
		onit add simulator -c my-cluster

		# Add a network of stratum switches that emulates a linear network topology with two nodes
		onit add network -n stratum-linear --topo linear --devices 2
	   
		# Add latest version of an application
		onit add app -n onos-ztp -i onosproject/onos-ztp:latest -u 0 -p grpc=5150 -r 2 -s /certs/onf.cacrt=configs/certs/onf.cacrt -s /certs/onos-ztp.crt=configs/certs/service.crt -s /certs/onos-ztp.key=configs/certs/service.key -- -caPath=/certs/onf.cacrt -keyPath=/certs/onos-ztp.key -certPath=/certs/onos-ztp.crt --image-pull-policy "Always" `
)

const (
	defaultMininetImage      = "opennetworking/mn-stratum:latest"
	defaultSimulatorImage    = "onosproject/device-simulator:latest"
	defaultRANSimulatorImage = "onosproject/ran-simulator:latest"
)

// getAddCommand returns a cobra "add" command for adding resources to the cluster
func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add {simulator,network} [args]",
		Short:   "Add resources to the cluster",
		Example: addExample,
	}
	cmd.AddCommand(getAddSimulatorCommand())
	cmd.AddCommand(getAddRANSimulatorCommand())
	cmd.AddCommand(getAddNetworkCommand())
	cmd.AddCommand(getAddAppCommand())
	return cmd
}

// getAddNetworkCommand returns a cobra command for deploying a stratum network
func getAddNetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network [args]",
		Short: "Add a stratum network to the test cluster",
		Args:  cobra.MaximumNArgs(10),
		RunE:  runInCluster(runAddNetworkCommand),
	}

	cmd.Flags().StringP("name", "n", "", "the name of the network to add")
	cmd.Flags().StringP("image", "i", defaultMininetImage, "the image to deploy")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringP("topo", "t", "", "the topology to create")
	_ = cmd.MarkFlagRequired("topo")
	cmd.Flags().IntP("devices", "d", 0, "the number of devices in the topology")
	_ = cmd.MarkFlagRequired("devices")
	return cmd
}

func runAddNetworkCommand(cmd *cobra.Command, _ []string) error {
	setupCommand(cmd)

	image, _ := cmd.Flags().GetString("image")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)
	topo, _ := cmd.Flags().GetString("topo")
	devices, _ := cmd.Flags().GetInt("devices")

	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	_, err = env.Networks().
		New().
		SetName(getName(cmd)).
		SetImage(image).
		SetPullPolicy(pullPolicy).
		SetCustom(topo, devices).
		Add()
	return err
}

// getAddSimulatorCommand returns a cobra command for deploying a device simulator
func getAddSimulatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulator [args]",
		Short: "Add a device simulator to the test cluster",
		Args:  cobra.NoArgs,
		RunE:  runInCluster(runAddSimulatorCommand),
	}

	cmd.Flags().StringP("name", "n", "", "the name to assign to the device simulator")
	cmd.Flags().StringP("image", "i", defaultSimulatorImage, "the image to deploy")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	return cmd
}

func runAddSimulatorCommand(cmd *cobra.Command, _ []string) error {
	setupCommand(cmd)

	image, _ := cmd.Flags().GetString("image")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	_, err = env.Simulators().
		New().
		SetName(getName(cmd)).
		SetPort(11161).
		SetImage(image).
		SetPullPolicy(pullPolicy).
		Add()
	return err
}

// getAddAppCommand returns a cobra command for deploying an app
func getAddAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app [args]",
		Short: "Add an app to the test cluster",
		Args:  cobra.ArbitraryArgs,
		RunE:  runInCluster(runAddAppCommand),
	}

	cmd.Flags().StringP("name", "n", "", "the name of the app to add")
	cmd.Flags().StringP("image", "i", "", "the image to deploy")
	_ = cmd.MarkFlagRequired("image")
	cmd.Flags().IntP("replicas", "r", 1, "the number of replicas to deploy")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringToIntP("port", "p", map[string]int{}, "ports to expose")
	cmd.Flags().BoolP("debug", "d", false, "enable debug mode")
	cmd.Flags().StringToStringP("secret", "s", map[string]string{}, "secrets to add to the application")
	cmd.Flags().Bool("privileged", false, "run the application in privileged mode")
	cmd.Flags().IntP("user", "u", -1, "set the user with which to run the application")
	cmd.Flags().StringToStringP("env", "e", map[string]string{}, "set environment variables")
	return cmd
}

func runAddAppCommand(cmd *cobra.Command, args []string) error {
	setupCommand(cmd)

	image, _ := cmd.Flags().GetString("image")
	replicas, _ := cmd.Flags().GetInt("replicas")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)
	ports, _ := cmd.Flags().GetStringToInt("port")
	debug, _ := cmd.Flags().GetBool("debug")
	privileged, _ := cmd.Flags().GetBool("privileged")
	user, _ := cmd.Flags().GetInt("user")
	envVars, _ := cmd.Flags().GetStringToString("env")
	secretValues, _ := cmd.Flags().GetStringToString("secret")

	secrets := make(map[string]string)
	for name, value := range secretValues {
		secret, err := ioutil.ReadFile(value)
		if err != nil {
			secrets[name] = value
		} else {
			secrets[name] = string(secret)
		}
	}

	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	setup := env.NewApp().
		SetName(getName(cmd)).
		SetReplicas(replicas).
		SetImage(image).
		SetPullPolicy(pullPolicy).
		SetPorts(ports).
		SetDebug(debug).
		SetSecrets(secrets).
		SetPrivileged(privileged).
		SetEnv(envVars).
		SetArgs(args...)
	if user >= 0 {
		setup.SetUser(user)
	}
	_, err = setup.Add()
	return err
}

// getAddRANSimulatorCommand returns a cobra command for deploying a device simulator
func getAddRANSimulatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ran-simulator [args]",
		Short: "Add a RAN simulator to the test cluster",
		Args:  cobra.NoArgs,
		RunE:  runInCluster(runAddRANSimulatorCommand),
	}

	cmd.Flags().StringP("name", "n", "", "the name to assign to the device simulator")
	cmd.Flags().StringP("image", "i", defaultRANSimulatorImage, "the image to deploy")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	return cmd
}

func runAddRANSimulatorCommand(cmd *cobra.Command, _ []string) error {
	setupCommand(cmd)

	image, _ := cmd.Flags().GetString("image")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	kubeAPI, err := kube.GetAPI(getCluster(cmd))
	if err != nil {
		return err
	}

	env := env.New(kubeAPI)
	_, err = env.Simulators().
		New().
		SetName("ran-simulator").
		SetPort(5150).
		SetImage(image).
		SetPullPolicy(pullPolicy).
		Add()
	return err
}
