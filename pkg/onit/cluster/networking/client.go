package networking

import (
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	"github.com/onosproject/onos-test/pkg/onit/cluster/networking/v1beta1"
)

type Client interface {
	NetworkingV1Beta1() v1beta1.Client
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

func (c *client) NetworkingV1Beta1() v1beta1.Client {
	return v1beta1.NewClient(c.ObjectsClient)
}
