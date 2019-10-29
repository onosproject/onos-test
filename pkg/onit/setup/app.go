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

	"github.com/onosproject/onos-test/pkg/onit/k8s"
	corev1 "k8s.io/api/core/v1"
)

// AddApp add an application to the cluster
func (t *TestSetup) AddApp() {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	pullPolicy := corev1.PullPolicy(t.imagePullPolicy)

	if pullPolicy != corev1.PullAlways && pullPolicy != corev1.PullIfNotPresent && pullPolicy != corev1.PullNever {
		exitError(fmt.Errorf("invalid pull policy; must of one of %s, %s or %s", corev1.PullAlways, corev1.PullIfNotPresent, corev1.PullNever))
	}

	// Create the app configuration
	config := &k8s.AppConfig{
		Image:      t.imageName,
		PullPolicy: pullPolicy,
	}

	// Add the app to the cluster
	if status := cluster.AddApp(t.appName, config); status.Failed() {
		exitStatus(status)
	} else {
		fmt.Println(t.appName)
	}
}

// RemoveApp remove an app from the cluster
func (t *TestSetup) RemoveApp() {
	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}

	apps, err := cluster.GetApps()
	if err != nil {
		exitError(err)
	}

	if !Contains(apps, t.appName) {
		exitError(errors.New("the given app name does not exist"))
	}

	// Remove the app from the cluster
	if status := cluster.RemoveApp(t.appName); status.Failed() {
		exitStatus(status)
	}

}

// GetApps return the list current apps in the cluster
func (t *TestSetup) GetApps() ([]string, error) {

	controller := t.initController()
	// Get the cluster controller
	cluster, err := controller.GetCluster(t.clusterID)
	if err != nil {
		exitError(err)
	}
	return cluster.GetApps()
}
