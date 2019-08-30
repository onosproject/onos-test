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
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// SetImage updates the container image of a deployment
func (c *ClusterController) SetImage(deploymentID string, image string, imagePullPolicy string) {

	deploymentsClient := c.kubeclient.AppsV1().Deployments(c.clusterID)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(deploymentID, metav1.GetOptions{})
		if getErr != nil {
			c.status.Fail(getErr)
		}
		result.Spec.Template.Spec.Containers[0].Image = image
		result.Spec.Template.Spec.Containers[0].ImagePullPolicy = v1.PullPolicy(imagePullPolicy)

		if strings.Compare(result.Spec.Template.Spec.Containers[0].Image, image) == 0 {
			result.Spec.Template.CreationTimestamp = metav1.Now()
		}
		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})

	if retryErr != nil {
		c.status.Fail(retryErr)
	}
}
