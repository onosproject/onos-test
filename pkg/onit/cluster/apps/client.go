package apps

import (
	"github.com/onosproject/onos-test/pkg/onit/cluster/apps/v1"
	"github.com/onosproject/onos-test/pkg/onit/cluster/apps/v1beta1"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
	AppsV1() v1.Client
	AppsV1Beta1() v1beta1.Client
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

func (c *client) AppsV1() v1.Client {
	return v1.NewClient(c.ObjectsClient)
}

func (c *client) AppsV1Beta1() v1beta1.Client {
	return v1beta1.NewClient(c.ObjectsClient)
}
