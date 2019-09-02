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

package k8s

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/onit/console"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"strings"
	"time"
)

// SetImage updates the container image of a resource
func (c *ClusterController) SetImage(resourceID string, image string, pullPolicy corev1.PullPolicy) console.ErrorStatus {
	c.status.Start("Loading pod configurations")
	pods, err := c.kubeclient.CoreV1().Pods(c.clusterID).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=onos,resource=%s", resourceID),
	})
	if err != nil {
		return c.status.Fail(err)
	} else if len(pods.Items) == 0 {
		return c.status.Fail(fmt.Errorf("no resources matching '%s' found", resourceID))
	}
	c.status.Succeed()

	// Attempt to determine whether this is a deployment by checking for a deployment with the same name as the resource ID
	deployment, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Get(resourceID, metav1.GetOptions{})

	// If no deployment was found, assume this is a pod or another type of pod set
	if err != nil && k8serrors.IsNotFound(err) {
		return c.setPodsImage(pods, image, pullPolicy)
	}
	return c.setDeploymentImage(deployment, pods, image, pullPolicy)
}

// setPodsImage updates the image for a set of pods
func (c *ClusterController) setPodsImage(pods *corev1.PodList, image string, pullPolicy corev1.PullPolicy) console.ErrorStatus {
	// First, update all the pods
	c.status.Start("Updating pods")
	for _, pod := range pods.Items {
		c.status.Progress(fmt.Sprintf("Updating %s", pod.Name))
		if err := c.updatePod(pod, image, pullPolicy); err != nil {
			return c.status.Fail(err)
		}
	}
	c.status.Succeed()

	c.status.Start("Waiting for pods to become ready")
	ready := make(map[string]bool)
	total := len(pods.Items)
	c.status.Progress(fmt.Sprintf("0/%d", total))

	// Loop through the updated pods and check their statuses until all pods are ready
	state := ""
	for len(ready) < total {
		stateUpdated := false
		for _, pod := range pods.Items {
			if _, ok := ready[pod.Name]; !ok {
				if podReady, err := c.podReady(pod.Name); err == nil {
					if podReady {
						ready[pod.Name] = true
						c.status.Progress(fmt.Sprintf("(%d/%d) %s", len(ready), total, state))
					} else if !stateUpdated {
						state = fmt.Sprintf("%s: %s", pod.Name, pod.Status.Phase)
						c.status.Progress(fmt.Sprintf("(%d/%d) %s", len(ready), total, state))
						stateUpdated = true
					}
				} else {
					return c.status.Fail(err)
				}
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	return c.status.Succeed()
}

// setDeploymentImage updates the image for a deployment
func (c *ClusterController) setDeploymentImage(deployment *appsv1.Deployment, pods *corev1.PodList, image string, pullPolicy corev1.PullPolicy) console.ErrorStatus {
	// If this is a deployment, update the deployment
	c.status.Start("Updating deployment")
	if err := c.updateDeployment(deployment.Name, image, pullPolicy); err != nil {
		return c.status.Fail(err)
	}
	c.status.Succeed()

	// Once the deployment has been updated, loop and block until all the pods in the deployment have been updated
	c.status.Start("Waiting for pods to become ready")
	ready := make(map[string]bool)
	total := int(*deployment.Spec.Replicas)
	state := ""
	c.status.Progress(fmt.Sprintf("0/%d", total))
	for len(ready) < total {
		// Get the list of pods matching the deployment
		updates, err := c.kubeclient.CoreV1().Pods(c.clusterID).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=onos,resource=%s", deployment.Name),
		})
		if err != nil {
			return c.status.Fail(err)
		}

		stateUpdated := false
		for _, pod := range updates.Items {
			// If the pod has not already been marked ready, check if all its containers are ready
			if _, ok := ready[pod.Name]; !ok {
				// If the pod is in the original pods list, ignore the pod.
				if c.podListContains(pods, pod) {
					continue
				}

				if podReady := c.containersReady(pod); podReady {
					ready[pod.Name] = true
					c.status.Progress(fmt.Sprintf("(%d/%d) %s", len(ready), total, state))
				} else if !stateUpdated {
					state = fmt.Sprintf("%s: %s", pod.Name, pod.Status.Phase)
					c.status.Progress(fmt.Sprintf("(%d/%d) %s", len(ready), total, state))
					stateUpdated = true
				}
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	return c.status.Succeed()
}

// updatePod updates the given pod
func (c *ClusterController) updatePod(pod corev1.Pod, image string, pullPolicy corev1.PullPolicy) error {
	if err := c.kubeclient.CoreV1().Pods(c.clusterID).Delete(pod.Name, &metav1.DeleteOptions{}); err != nil {
		return err
	}

	update := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        pod.Name,
			Namespace:   pod.Namespace,
			Labels:      pod.Labels,
			Annotations: pod.Annotations,
		},
		Spec: pod.Spec,
	}
	update.Spec.Containers[0].Image = image
	update.Spec.Containers[0].ImagePullPolicy = pullPolicy
	for {
		_, err := c.kubeclient.CoreV1().Pods(c.clusterID).Create(update)
		if err == nil || !k8serrors.IsAlreadyExists(err) {
			return err
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// podReady checks whether the pod of the given name is ready
func (c *ClusterController) podReady(name string) (bool, error) {
	// Get the current state of the pod
	pod, err := c.kubeclient.CoreV1().Pods(c.clusterID).Get(name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return c.containersReady(*pod), nil
}

// containersReady returns a bool indicating whether all the given pod's containers are ready
func (c *ClusterController) containersReady(pod corev1.Pod) bool {
	if len(pod.Status.ContainerStatuses) != len(pod.Spec.Containers) {
		return false
	}
	for _, status := range pod.Status.ContainerStatuses {
		if !status.Ready {
			return false
		}
	}
	return true
}

// updateDeployment updates the given deployment
func (c *ClusterController) updateDeployment(name string, image string, pullPolicy corev1.PullPolicy) error {
	// Update the deployment image and modify the creation timestamp to ensure the update is applied regardless
	// of whether the image name has changed
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deployment, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		deployment.Spec.Template.Spec.Containers[0].Image = image
		deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy = corev1.PullPolicy(pullPolicy)

		if strings.Compare(deployment.Spec.Template.Spec.Containers[0].Image, image) == 0 {
			deployment.Spec.Template.CreationTimestamp = metav1.Now()
		}
		_, updateErr := c.kubeclient.AppsV1().Deployments(c.clusterID).Update(deployment)
		return updateErr
	})
}

// podListContains returns a bool indicating whether the given pod is contained in the given list of pods
func (c *ClusterController) podListContains(list *corev1.PodList, pod corev1.Pod) bool {
	for _, item := range list.Items {
		if pod.Name == item.Name {
			return true
		}
	}
	return false
}
