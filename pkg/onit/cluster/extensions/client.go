package extensions

import (
    extensionsv1beta1 "github.com/onosproject/onos-test/pkg/onit/cluster/extensions/v1beta1"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
    ExtensionsV1Beta1() extensionsv1beta1.Client
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}


func (c *client) ExtensionsV1Beta1() extensionsv1beta1.Client {
	return extensionsv1beta1.NewClient(c.ObjectsClient)
}

