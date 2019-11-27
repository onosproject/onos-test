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
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/onosproject/onos-test/pkg/onit/console"

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
	switch k8s.NodeType(t.nodeType) {
	case k8s.OnosAll:
		return cluster.GetNodes()
	case k8s.OnosConfig:
		return cluster.GetOnosConfigNodes()
	case k8s.OnosTopo:
		return cluster.GetOnosTopoNodes()
	case k8s.OnosCli:
		return cluster.GetOnosCliNodes()
	case k8s.OnosGui:
		return cluster.GetOnosGuiNodes()
	}

	return nil, errors.New("Unsupported node type " + t.nodeType)
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

// OpenDebug open a debug session for a given node
func (t *TestSetup) OpenDebug() {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	if len(t.args) == 0 {
		var wg sync.WaitGroup
		nodes, _ := t.GetNodes()
		n := len(nodes)
		wg.Add(n)

		asyncErrors := make(chan error, n)
		freePorts, err := k8s.GetFreePorts(n)
		if err != nil {
			exitError(err)
		}

		for index := range nodes {
			go func(node k8s.NodeInfo, port int) {
				fmt.Println("Start Debugger for:", node.ID)
				err = cluster.PortForward(node.ID, port, 40000)
				asyncErrors <- err
				wg.Done()
			}(nodes[index], freePorts[index])

		}

		go func() {
			wg.Wait()
			close(asyncErrors)
		}()

		for err = range asyncErrors {
			if err != nil {
				exitError(err)
			}
		}

	} else {
		if err := cluster.PortForward(t.args[0], t.debugPort, 40000); err != nil {
			exitError(err)
		} else {
			fmt.Println(t.debugPort)
		}

	}
}

// FetchLogs fetch logs for a given resourceID and writes them to a log file
func (t *TestSetup) FetchLogs() {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	if len(t.args) > 0 {
		resourceID := t.args[0]
		resources, err := cluster.GetResources(resourceID)
		if err != nil {
			exitError(err)
		}

		var status console.ErrorStatus
		for _, resource := range resources {
			path := filepath.Join(t.logDestination, fmt.Sprintf("%s.log", resource))
			status = cluster.DownloadLogs(resource, path, t.logOptions)
		}

		if status.Failed() {
			exitStatus(status)
		}
	} else {
		nodes, err := cluster.GetNodes()
		if err != nil {
			exitError(err)
		}

		var status console.ErrorStatus
		for _, node := range nodes {
			path := filepath.Join(t.logDestination, fmt.Sprintf("%s.log", node.ID))
			status = cluster.DownloadLogs(node.ID, path, t.logOptions)
		}

		if status.Failed() {
			exitStatus(status)
		}
	}

}
