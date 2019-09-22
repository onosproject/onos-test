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

package setup

import (
	"fmt"

	interfaces "github.com/onosproject/onos-test/pkg/onit/controller"
	"github.com/onosproject/onos-test/pkg/onit/k8s"
	"k8s.io/client-go/rest"
)

// TestSetup a struct to store test setup info
type TestSetup struct {
	clusterID       string
	dockerRegistry  string
	configNodes     int
	topoNodes       int
	partitions      int
	partitionSize   int
	configName      string
	imageTags       map[string]string
	imagePullPolicy string
	clusterType     string
	simulatorName   string
	appName         string
	imageName       string
	mininetOptions  []string
	networkName     string
}

// TestSetupBuilder test setup builder interface
type TestSetupBuilder interface {
	SetClusterID(string) TestSetupBuilder
	SetDockerRegistry(string) TestSetupBuilder
	SetConfigNodes(int) TestSetupBuilder
	SetTopoNodes(int) TestSetupBuilder
	SetPartitions(int) TestSetupBuilder
	SetPartitionSize(int) TestSetupBuilder
	SetConfigName(string) TestSetupBuilder
	SetImageTags(map[string]string) TestSetupBuilder
	SetImagePullPolicy(string) TestSetupBuilder
	SetClusterType(string) TestSetupBuilder
	SetSimulatorName(string) TestSetupBuilder
	SetAppName(string) TestSetupBuilder
	SetImageName(string) TestSetupBuilder
	SetMininetOptions([]string) TestSetupBuilder
	SetNetworkName(string) TestSetupBuilder
	Build() TestSetup
}

// SetClusterID set cluster ID
func (t *TestSetup) SetClusterID(clusterID string) TestSetupBuilder {

	t.clusterID = clusterID
	return t
}

// SetDockerRegistry set docker registry
func (t *TestSetup) SetDockerRegistry(dockerRegistery string) TestSetupBuilder {

	t.dockerRegistry = dockerRegistery
	return t
}

// SetConfigNodes set number of config nodes
func (t *TestSetup) SetConfigNodes(configNodes int) TestSetupBuilder {

	t.configNodes = configNodes
	return t
}

// SetTopoNodes set number of topo nodes
func (t *TestSetup) SetTopoNodes(topoNodes int) TestSetupBuilder {

	t.topoNodes = topoNodes
	return t
}

// SetPartitions set number of partitions
func (t *TestSetup) SetPartitions(partitions int) TestSetupBuilder {

	t.partitions = partitions
	return t
}

// SetConfigName set config name
func (t *TestSetup) SetConfigName(configName string) TestSetupBuilder {

	t.configName = configName
	return t
}

// SetImageTags set image tags
func (t *TestSetup) SetImageTags(imageTags map[string]string) TestSetupBuilder {

	t.imageTags = imageTags
	return t
}

// SetClusterType set cluster type
func (t *TestSetup) SetClusterType(clusteType string) TestSetupBuilder {

	t.clusterType = clusteType
	return t
}

// SetImagePullPolicy set image pull policy
func (t *TestSetup) SetImagePullPolicy(imagePullPolicy string) TestSetupBuilder {

	t.imagePullPolicy = imagePullPolicy
	return t
}

// SetPartitionSize set size of the partition
func (t *TestSetup) SetPartitionSize(partitionSize int) TestSetupBuilder {

	t.partitionSize = partitionSize
	return t
}

// SetSimulatorName set the name of the simulator
func (t *TestSetup) SetSimulatorName(name string) TestSetupBuilder {
	t.simulatorName = name
	return t
}

// SetAppName set an application name
func (t *TestSetup) SetAppName(appName string) TestSetupBuilder {
	t.appName = appName
	return t

}

// SetImageName set the name of app image
func (t *TestSetup) SetImageName(imageName string) TestSetupBuilder {
	t.imageName = imageName
	return t
}

// SetMininetOptions set mininet options for a stratum network
func (t *TestSetup) SetMininetOptions(mininetOptions []string) TestSetupBuilder {
	t.mininetOptions = mininetOptions
	return t
}

// SetNetworkName set a name for a stratum network
func (t *TestSetup) SetNetworkName(networkName string) TestSetupBuilder {
	t.networkName = networkName
	return t
}

// New creates an instance of TestSetupBuilder with default values
func New() TestSetupBuilder {
	imageTags := make(map[string]string)
	InitImageTags(imageTags)
	clusterID := fmt.Sprintf("cluster-%s", NewUUIDString())
	return &TestSetup{
		clusterID:       clusterID,
		dockerRegistry:  "",
		configNodes:     1,
		topoNodes:       1,
		partitions:      1,
		partitionSize:   1,
		configName:      "default",
		imagePullPolicy: "IfNotPresent",
		clusterType:     "k8s",
		imageTags:       imageTags,
	}
}

// Build build an instance of testSetup
func (t *TestSetup) Build() TestSetup {
	return TestSetup{
		clusterID:       t.clusterID,
		dockerRegistry:  t.dockerRegistry,
		configNodes:     t.configNodes,
		topoNodes:       t.topoNodes,
		partitions:      t.partitions,
		partitionSize:   t.partitionSize,
		configName:      t.configName,
		imagePullPolicy: t.imagePullPolicy,
		clusterType:     t.clusterType,
		imageTags:       t.imageTags,
		simulatorName:   t.simulatorName,
		imageName:       t.imageName,
		networkName:     t.networkName,
		appName:         t.appName,
		mininetOptions:  t.mininetOptions,
	}
}

// initController creates an instance of controller interface and initialize it
func (t *TestSetup) initController() interfaces.Controller {
	var controller interfaces.Controller
	if t.clusterType == string(k8s.K8s) {
		k8sController, err := k8s.NewController()
		if err != nil {
			exitError(err)
		}
		controller = k8sController
	}
	return controller
}

// GetRestConfig returns the k8s config
func (t *TestSetup) GetRestConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}
