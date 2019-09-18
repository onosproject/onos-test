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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

// InitImageTags initialize the default values of image tags
func InitImageTags(imageTags map[string]string) {
	if imageTags["config"] == "" {
		imageTags["config"] = string(k8s.Debug)
	}
	if imageTags["topo"] == "" {
		imageTags["topo"] = string(k8s.Debug)
	}
	if imageTags["gui"] == "" {
		imageTags["gui"] = string(k8s.Latest)
	}
	if imageTags["cli"] == "" {
		imageTags["cli"] = string(k8s.Latest)
	}
	if imageTags["atomix"] == "" {
		imageTags["atomix"] = string(k8s.Latest)
	}
	if imageTags["raft"] == "" {
		imageTags["raft"] = string(k8s.Latest)
	}
	if imageTags["simulator"] == "" {
		imageTags["simulator"] = string(k8s.Latest)
	}
	if imageTags["stratum"] == "" {
		imageTags["stratum"] = string(k8s.Latest)
	}
	if imageTags["test"] == "" {
		imageTags["test"] = string(k8s.Latest)
	}

}

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
		configName:      "config",
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
	}
}

// CreateCluster creates a k8s cluster
func (t *TestSetup) CreateCluster() error {
	var controller interfaces.Controller
	if t.clusterType == string(k8s.K8s) {
		k8sController, err := k8s.NewController()
		if err != nil {
			exitError(err)
		}
		controller = k8sController
	}

	pullPolicy := corev1.PullPolicy(t.imagePullPolicy)

	if pullPolicy != corev1.PullAlways && pullPolicy != corev1.PullIfNotPresent && pullPolicy != corev1.PullNever {
		exitError(fmt.Errorf("invalid pull policy; must of one of %s, %s or %s", corev1.PullAlways, corev1.PullIfNotPresent, corev1.PullNever))
	}

	// Create the cluster configuration
	config := &k8s.ClusterConfig{
		Registry:      t.dockerRegistry,
		Preset:        t.configName,
		ImageTags:     t.imageTags,
		PullPolicy:    pullPolicy,
		ConfigNodes:   t.configNodes,
		TopoNodes:     t.topoNodes,
		Partitions:    t.partitions,
		PartitionSize: t.partitionSize,
	}

	// Create the cluster controller
	cluster, status := controller.NewCluster(t.clusterID, config)
	if status.Failed() {
		exitStatus(status)
	}

	var k8sCluster interfaces.ClusterController = cluster

	// Setup the cluster
	if status := k8sCluster.Setup(); status.Failed() {
		exitStatus(status)
	} else {
		fmt.Println(t.clusterID)
	}

	return nil
}

// GetRestConfig returns the k8s config
func (t *TestSetup) GetRestConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}
