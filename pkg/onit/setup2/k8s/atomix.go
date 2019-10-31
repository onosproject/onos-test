package k8s

import "time"

// setupAtomixController sets up the Atomix controller and associated resources
func (c *TestSetup) setupAtomixController() error {
	if err := c.createAtomixPartitionSetResource(); err != nil {
		return err
	}
	if err := c.createAtomixPartitionResource(); err != nil {
		return err
	}
	if err := c.createAtomixDeployment(); err != nil {
		return err
	}
	if err := c.createAtomixService(); err != nil {
		return err
	}
	if err := c.awaitAtomixControllerReady(); err != nil {
		return err
	}
	return nil
}

// createAtomixPartitionSetResource creates the PartitionSet custom resource definition in the k8s cluster
func (c *ClusterController) createAtomixPartitionSetResource() error {
	crd := &apiextensionv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "partitionsets.k8s.atomix.io",
		},
		Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
			Group: "k8s.atomix.io",
			Names: apiextensionv1beta1.CustomResourceDefinitionNames{
				Kind:     "PartitionSet",
				ListKind: "PartitionSetList",
				Plural:   "partitionsets",
				Singular: "partitionset",
			},
			Scope:   apiextensionv1beta1.NamespaceScoped,
			Version: "v1alpha1",
			Subresources: &apiextensionv1beta1.CustomResourceSubresources{
				Status: &apiextensionv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}

	_, err := c.extensionsclient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createAtomixPartitionResource creates the Partition custom resource definition in the k8s cluster
func (c *ClusterController) createAtomixPartitionResource() error {
	crd := &apiextensionv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "partitions.k8s.atomix.io",
		},
		Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
			Group: "k8s.atomix.io",
			Names: apiextensionv1beta1.CustomResourceDefinitionNames{
				Kind:     "Partition",
				ListKind: "PartitionList",
				Plural:   "partitions",
				Singular: "partition",
			},
			Scope:   apiextensionv1beta1.NamespaceScoped,
			Version: "v1alpha1",
			Subresources: &apiextensionv1beta1.CustomResourceSubresources{
				Status: &apiextensionv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}

	_, err := c.extensionsclient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createAtomixDeployment creates the Atomix controller Deployment
func (c *ClusterController) createAtomixDeployment() error {
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "atomix-controller",
			Namespace: c.clusterID,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "atomix-controller",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": "atomix-controller",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: c.clusterID,
					Containers: []corev1.Container{
						{
							Name:            "atomix-controller",
							Image:           c.imageName("atomix/atomix-k8s-controller", c.config.ImageTags["atomix"]),
							ImagePullPolicy: c.config.PullPolicy,
							Command:         []string{"atomix-controller"},
							Env: []corev1.EnvVar{
								{
									Name:  "CONTROLLER_NAME",
									Value: "atomix-controller",
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
									ContainerPort: 5679,
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
								InitialDelaySeconds: int32(4),
								PeriodSeconds:       int32(10),
								FailureThreshold:    int32(1),
							},
						},
					},
				},
			},
		},
	}
	_, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Create(deployment)
	return err
}

// createAtomixService creates a service for the controller
func (c *ClusterController) createAtomixService() error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "atomix-controller",
			Namespace: c.clusterID,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"name": "atomix-controller",
			},
			Ports: []corev1.ServicePort{
				{
					Name: "control",
					Port: 5679,
				},
			},
		},
	}
	_, err := c.kubeclient.CoreV1().Services(c.clusterID).Create(service)
	return err
}

// awaitAtomixControllerReady blocks until the Atomix controller is ready
func (c *ClusterController) awaitAtomixControllerReady() error {
	for {
		dep, err := c.kubeclient.AppsV1().Deployments(c.clusterID).Get("atomix-controller", metav1.GetOptions{})
		if err != nil {
			return err
		} else if dep.Status.ReadyReplicas == 1 {
			return nil
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
