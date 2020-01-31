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
	"time"

	"github.com/onosproject/onos-test/pkg/util/logging"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	atomixType    = "atomix"
	atomixImage   = "atomix/kubernetes-controller:latest"
	atomixService = "atomix-controller"
	atomixPort    = 5679
)

func newAtomix(cluster *Cluster) *Atomix {
	deployment := newDeployment(cluster)
	deployment.SetLabels(getLabels(atomixType))
	deployment.SetImage(atomixImage)
	deployment.SetName(atomixService)

	return &Atomix{
		Deployment: deployment,
	}
}

// Atomix provides methods for managing the Atomix controller
type Atomix struct {
	*Deployment
}

// Setup sets up the Atomix controller and associated resources
func (s *Atomix) Setup() error {
	step := logging.NewStep(s.namespace, "Setup Atomix controller")
	step.Start()
	step.Log("Creating Database resource")
	if err := s.createDatabaseResource(); err != nil {
		step.Fail(err)
		return err
	}
	step.Log("Creating Cluster resource")
	if err := s.createClusterResource(); err != nil {
		step.Fail(err)
		return err
	}
	step.Log("Creating Partition resource")
	if err := s.createPartitionResource(); err != nil {
		step.Fail(err)
		return err
	}
	step.Log("Creating controller Deployment")
	if err := s.createDeployment(); err != nil {
		step.Fail(err)
		return err
	}
	step.Log("Creating controller Service")
	if err := s.createService(); err != nil {
		step.Fail(err)
		return err
	}
	step.Log("Waiting for controller to become ready")
	if err := s.awaitReady(); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// createDatabaseResource creates the PartitionSet custom resource definition in the k8s cluster
func (s *Atomix) createDatabaseResource() error {
	crd := &apiextensionv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "databases.cloud.atomix.io",
		},
		Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
			Group: "cloud.atomix.io",
			Names: apiextensionv1beta1.CustomResourceDefinitionNames{
				Kind:     "Database",
				ListKind: "DatabaseList",
				Plural:   "databases",
				Singular: "database",
			},
			Scope:   apiextensionv1beta1.NamespaceScoped,
			Version: "v1beta1",
			Subresources: &apiextensionv1beta1.CustomResourceSubresources{
				Status: &apiextensionv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}

	_, err := s.extensionsClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createClusterResource creates the Cluster custom resource definition in the k8s cluster
func (s *Atomix) createClusterResource() error {
	crd := &apiextensionv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "clusters.cloud.atomix.io",
		},
		Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
			Group: "cloud.atomix.io",
			Names: apiextensionv1beta1.CustomResourceDefinitionNames{
				Kind:     "Cluster",
				ListKind: "ClusterList",
				Plural:   "clusters",
				Singular: "cluster",
			},
			Scope:   apiextensionv1beta1.NamespaceScoped,
			Version: "v1beta1",
			Subresources: &apiextensionv1beta1.CustomResourceSubresources{
				Status: &apiextensionv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}

	_, err := s.extensionsClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createPartitionResource creates the Partition custom resource definition in the k8s cluster
func (s *Atomix) createPartitionResource() error {
	crd := &apiextensionv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "partitions.cloud.atomix.io",
		},
		Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
			Group: "cloud.atomix.io",
			Names: apiextensionv1beta1.CustomResourceDefinitionNames{
				Kind:     "Partition",
				ListKind: "PartitionList",
				Plural:   "partitions",
				Singular: "partition",
			},
			Scope:   apiextensionv1beta1.NamespaceScoped,
			Version: "v1beta1",
			Subresources: &apiextensionv1beta1.CustomResourceSubresources{
				Status: &apiextensionv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}

	_, err := s.extensionsClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createDeployment creates the Atomix controller Deployment
func (s *Atomix) createDeployment() error {
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name(),
			Namespace: s.namespace,
			Labels:    s.labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: s.labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: s.labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: s.namespace,
					Containers: []corev1.Container{
						{
							Name:            s.Name(),
							Image:           s.Image(),
							ImagePullPolicy: s.PullPolicy(),
							Command:         []string{"atomix-controller"},
							Env: []corev1.EnvVar{
								{
									Name:  "CONTROLLER_NAME",
									Value: s.name,
								},
								{
									Name: "CONTROLLER_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "control",
									ContainerPort: atomixPort,
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"stat",
											"/tmp/atomix-controller-ready",
										},
									},
								},
								PeriodSeconds:    1,
								FailureThreshold: 30,
							},
						},
					},
				},
			},
		},
	}
	_, err := s.kubeClient.AppsV1().Deployments(s.namespace).Create(deployment)
	return err
}

// createService creates a service for the controller
func (s *Atomix) createService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name(),
			Namespace: s.namespace,
			Labels:    s.labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: s.labels,
			Ports: []corev1.ServicePort{
				{
					Name: "control",
					Port: atomixPort,
				},
			},
		},
	}
	_, err := s.kubeClient.CoreV1().Services(s.namespace).Create(service)
	return err
}

// awaitReady blocks until the Atomix controller is ready
func (s *Atomix) awaitReady() error {
	for {
		dep, err := s.kubeClient.AppsV1().Deployments(s.namespace).Get(s.name, metav1.GetOptions{})
		if err != nil {
			return err
		} else if dep.Status.ReadyReplicas == 1 {
			return nil
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
