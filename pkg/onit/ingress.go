package onit

import (
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// setupIngress sets up the Ingress
func (c *ClusterController) setupIngress() error {
	return c.createIngress()
}

// createIngress creates an ingress for onos services
func (c *ClusterController) createIngress() error {
	ing := &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "onos-ingress",
			Namespace: c.clusterID,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
				"nginx.org/grpc-services": "onos-config,onos-topo",
				"nginx.ingress.kubernetes.io/backend-protocol": "GRPCS",
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