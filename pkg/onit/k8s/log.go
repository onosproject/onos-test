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
	"bytes"
	"io"
	"os"

	"github.com/onosproject/onos-test/pkg/onit/console"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetLogs returns the logs for a single test resource
func (c *ClusterController) GetLogs(resourceID string, options corev1.PodLogOptions) ([]byte, error) {
	pod, err := c.Kubeclient.CoreV1().Pods(c.ClusterID).Get(resourceID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return c.getLogs(*pod, options)
}

// getLogs gets the logs from the given pod
func (c *ClusterController) getLogs(pod corev1.Pod, options corev1.PodLogOptions) ([]byte, error) {
	req := c.Kubeclient.CoreV1().Pods(c.ClusterID).GetLogs(pod.Name, &options)
	readCloser, err := req.Stream()
	if err != nil {
		return nil, err
	}

	defer readCloser.Close()

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(readCloser); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// StreamLogs streams the logs for the given test resources to stdout
func (c *ClusterController) StreamLogs(resourceID string) (io.ReadCloser, error) {
	pod, err := c.Kubeclient.CoreV1().Pods(c.ClusterID).Get(resourceID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return c.streamLogs(*pod)
}

// streamLogs streams the logs from the given pod to stdout
func (c *ClusterController) streamLogs(pod corev1.Pod) (io.ReadCloser, error) {
	req := c.Kubeclient.CoreV1().Pods(c.ClusterID).GetLogs(pod.Name, &corev1.PodLogOptions{
		Follow: true,
	})
	return req.Stream()
}

// DownloadLogs downloads the logs for the given resource to the given path
func (c *ClusterController) DownloadLogs(resourceID string, path string, options corev1.PodLogOptions) console.ErrorStatus {
	c.Status.Start("Downloading logs")
	pod, err := c.Kubeclient.CoreV1().Pods(c.ClusterID).Get(resourceID, metav1.GetOptions{})
	if err != nil {
		return c.Status.Fail(err)
	}
	if err := c.downloadLogs(*pod, path, options); err != nil {
		return c.Status.Fail(err)
	}
	return c.Status.Succeed()
}

// downloadLogs downloads the logs from the given pod to the given path
func (c *ClusterController) downloadLogs(pod corev1.Pod, path string, options corev1.PodLogOptions) error {
	// Create the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	// Get a stream of logs
	req := c.Kubeclient.CoreV1().Pods(c.ClusterID).GetLogs(pod.Name, &options)
	readCloser, err := req.Stream()
	if err != nil {
		return err
	}

	defer readCloser.Close()

	_, err = io.Copy(file, readCloser)
	return err
}
