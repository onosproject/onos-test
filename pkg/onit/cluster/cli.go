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

package cluster

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/util/logging"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	cliType    = "cli"
	cliImage   = "onosproject/onos-cli:latest"
	cliService = "onos-cli"
)

func newCLI(client *client) *CLI {
	return &CLI{
		Deployment: newDeployment(cliService, getLabels(cliType), cliImage, client),
	}
}

// CLI provides methods for managing the onos-cli service
type CLI struct {
	*Deployment
	enabled bool
}

// Enabled indicates whether the CLI is enabled
func (c *CLI) Enabled() bool {
	return c.enabled
}

// SetEnabled sets whether the CLI is enabled
func (c *CLI) SetEnabled(enabled bool) {
	c.enabled = enabled
}

// Setup sets up the CLI subsystem
func (c *CLI) Setup() error {
	if !c.enabled {
		return nil
	}

	step := logging.NewStep(c.namespace, "Setup onos-cli service")
	step.Start()
	step.Log("Creating onos-cli ConfigMap")
	if err := c.createConfigMap(); err != nil {
		step.Fail(err)
		return err
	}
	step.Log("Creating onos-cli Deployment")
	if err := c.createDeployment(); err != nil {
		step.Fail(err)
		return err
	}
	step.Log("Waiting for onos-cli to become ready")
	if err := c.AwaitReady(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// createConfigMap creates a ConfigMap to configure the onos-cli service
func (c *CLI) createConfigMap() error {
	config := fmt.Sprintf(`
controller: atomix-controller.%s.svc.cluster.local:5679
namespace: %s
group: raft
app: default
`, c.namespace, c.namespace)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cliService,
			Namespace: c.namespace,
			Labels:    getLabels(cliType),
		},
		Data: map[string]string{
			"atomix.yaml": config,
		},
	}
	_, err := c.kubeClient.CoreV1().ConfigMaps(c.namespace).Create(cm)
	return err
}

// createDeployment creates an onos-topo Deployment
func (c *CLI) createDeployment() error {
	nodes := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cliService,
			Namespace: c.namespace,
			Labels:    getLabels(cliType),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nodes,
			Selector: &metav1.LabelSelector{
				MatchLabels: getLabels(cliType),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getLabels(cliType),
					Annotations: map[string]string{
						"seccomp.security.alpha.kubernetes.io/pod": "unconfined",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            cliService,
							Image:           c.Image(),
							ImagePullPolicy: c.PullPolicy(),
							Stdin:           true,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/home/onos/.atomix/config.yaml",
									SubPath:   "atomix.yaml",
									ReadOnly:  false,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: cliService,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := c.kubeClient.AppsV1().Deployments(c.namespace).Create(deployment)
	return err
}

// AwaitReady waits for the onos-cli pods to complete startup
func (c *CLI) AwaitReady() error {
	if !c.enabled {
		return nil
	}

	for {
		// Get the onos-topo deployment
		dep, err := c.kubeClient.AppsV1().Deployments(c.namespace).Get(cliService, metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Return once the all replicas in the deployment are ready
		if dep.Status.ReadyReplicas > 0 {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}
