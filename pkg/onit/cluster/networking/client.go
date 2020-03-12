package networking

import (
    networkingv1beta1 "github.com/onosproject/onos-test/pkg/onit/cluster/networking/v1beta1"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
    NetworkingV1Beta1() networkingv1beta1.Client
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}


func (c *client) NetworkingV1Beta1() networkingv1beta1.Client {
	return networkingv1beta1.NewClient(c.ObjectsClient)
}

