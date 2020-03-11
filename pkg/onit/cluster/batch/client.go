package batch

import (
	"github.com/onosproject/onos-test/pkg/onit/cluster/batch/v1"
	"github.com/onosproject/onos-test/pkg/onit/cluster/batch/v1beta1"
	"github.com/onosproject/onos-test/pkg/onit/cluster/batch/v2alpha1"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
	BatchV1() v1.Client
	BatchV1Beta1() v1beta1.Client
	BatchV2Alpha1() v2alpha1.Client
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

func (c *client) BatchV1() v1.Client {
	return v1.NewClient(c.ObjectsClient)
}

func (c *client) BatchV1Beta1() v1beta1.Client {
	return v1beta1.NewClient(c.ObjectsClient)
}

func (c *client) BatchV2Alpha1() v2alpha1.Client {
	return v2alpha1.NewClient(c.ObjectsClient)
}
