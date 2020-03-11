package core

import (
	"github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
	CoreV1() v1.Client
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

// client is an implementation of the client interface
type client struct {
	metav1.ObjectsClient
}

func (c *client) CoreV1() v1.Client {
	return v1.NewClient(c.ObjectsClient)
}
