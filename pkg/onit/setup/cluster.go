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
	"strings"

	interfaces "github.com/onosproject/onos-test/pkg/onit/controller"
	"github.com/onosproject/onos-test/pkg/onit/k8s"
	corev1 "k8s.io/api/core/v1"
)

// CreateCluster creates a k8s cluster
func (t *TestSetup) CreateCluster() error {
	controller := t.initController()
	pullPolicy := corev1.PullPolicy(t.imagePullPolicy)

	if pullPolicy != corev1.PullAlways && pullPolicy != corev1.PullIfNotPresent && pullPolicy != corev1.PullNever {
		exitError(fmt.Errorf("invalid pull policy; must of one of %s, %s or %s", corev1.PullAlways, corev1.PullIfNotPresent, corev1.PullNever))
	}

	InitImageTags(t.imageTags)

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
	err := SetDefaultCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	// Setup the cluster
	if status := k8sCluster.Setup(); status.Failed() {
		exitStatus(status)
	} else {
		fmt.Println(t.clusterID)
	}

	return nil
}

// GetClusters returns the list of current clusters in the system
func (t *TestSetup) GetClusters() (map[string]*k8s.ClusterConfig, error) {
	controller := t.initController()
	return controller.GetClusters()
}

// GetCluster returns the current active cluster in the system
func (t *TestSetup) GetCluster() (*k8s.ClusterController, error) {
	controller := t.initController()
	return controller.GetCluster(t.clusterID)
}

// DeleteCluster delete the current cluster
func (t *TestSetup) DeleteCluster() {

	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)

	if err != nil {
		exitError(err)
	}

	status := cluster.Teardown()
	if status.Failed() {
		fmt.Println(status)
	} else {
		if err := SetDefaultCluster(""); err != nil {
			exitError(err)
		} else {
			fmt.Println(t.clusterID)
		}

	}
}

// SetCluster set the current clusterID to the given clusterID
func (t *TestSetup) SetCluster() {
	controller := t.initController()
	// Get the cluster controller
	_, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	// If we've made it this far, update the default cluster
	if err := SetDefaultCluster(t.clusterID); err != nil {
		exitError(err)
	} else {
		fmt.Println(t.clusterID)
	}
}

// GetNodes return the list of nodes in a cluster
func (t *TestSetup) GetNodes() ([]k8s.NodeInfo, error) {

	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	// Get the list of nodes based on given type
	if strings.Compare(t.nodeType, string(k8s.OnosAll)) == 0 {
		return cluster.GetNodes()

	} else if strings.Compare(t.nodeType, string(k8s.OnosConfig)) == 0 {
		return cluster.GetOnosConfigNodes()

	} else if strings.Compare(t.nodeType, string(k8s.OnosTopo)) == 0 {
		return cluster.GetOnosTopoNodes()

	} else if strings.Compare(t.nodeType, string(k8s.OnosCli)) == 0 {
		return cluster.GetOnosCliNodes()

	} else if strings.Compare(t.nodeType, string(k8s.OnosGui)) == 0 {
		return cluster.GetOnosGuiNodes()
	}

	return nil, nil

}

// GetHistory return a test history
func (t *TestSetup) GetHistory() ([]k8s.TestRecord, error) {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}
	return cluster.GetHistory()
}

// OpenSSH open a ssh session to a node for executing the remote commands
func (t *TestSetup) OpenSSH() {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}
	err = cluster.OpenShell(t.args[0], t.args[1:]...)
	if err != nil {
		exitError(err)
	}
}
