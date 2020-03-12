package core

import (
    corev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
    CoreV1() corev1.Client
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}


func (c *client) CoreV1() corev1.Client {
	return corev1.NewClient(c.ObjectsClient)
}

