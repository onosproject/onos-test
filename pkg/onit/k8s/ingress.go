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
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// setupIngress sets up the Ingress
func (c *ClusterController) setupIngress() error {
	if err := c.createGRPCIngress(); err != nil {
		return err
	}
	if err := c.createGUIIngress(); err != nil {
		return err
	}
	return nil
}

// createGRPCIngress creates an ingress for onos services
func (c *ClusterController) createGRPCIngress() error {
	ing := &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-ingress",
			Namespace: c.clusterID,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
				// gRPC services that can be routed by path
				"nginx.org/grpc-services": "onos-config,onos-topo",
				// Insecure backend gRPC protocol
				"nginx.ingress.kubernetes.io/backend-protocol": "GRPC",
			},
		},
		Spec: extensionsv1beta1.IngressSpec{
			TLS: []extensionsv1beta1.IngressTLS{
				{
					SecretName: c.clusterID,
				},
			},
			Rules: []extensionsv1beta1.IngressRule{
				{
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Path: "/gnmi.gNMI",
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: "onos-config",
										ServicePort: intstr.FromString("grpc"),
									},
								},
							},
						},
					},
				},
				{
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Path: "/proto.DeviceInventoryService",
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: "onos-config",
										ServicePort: intstr.FromString("grpc"),
									},
								},
							},
						},
					},
				},
				{
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Path: "/proto.DeviceService",
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: "onos-topo",
										ServicePort: intstr.FromString("grpc"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := c.kubeclient.ExtensionsV1beta1().Ingresses(c.clusterID).Create(ing)
	return err
}

// createGUIIngress creates an ingress for the GUI
func (c *ClusterController) createGUIIngress() error {
	ing := &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-gui-ingress",
			Namespace: c.clusterID,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
			},
		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{
				{
					Host: "onos-gui",
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: "onos-gui",
										ServicePort: intstr.FromInt(80),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := c.kubeclient.ExtensionsV1beta1().Ingresses(c.clusterID).Create(ing)
	return err
}
