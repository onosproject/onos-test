package batch

import (
    batchv1 "github.com/onosproject/onos-test/pkg/onit/cluster/batch/v1"
    batchv1beta1 "github.com/onosproject/onos-test/pkg/onit/cluster/batch/v1beta1"
    batchv2alpha1 "github.com/onosproject/onos-test/pkg/onit/cluster/batch/v2alpha1"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
    BatchV1() batchv1.Client
    BatchV1Beta1() batchv1beta1.Client
    BatchV2Alpha1() batchv2alpha1.Client
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}


func (c *client) BatchV1() batchv1.Client {
	return batchv1.NewClient(c.ObjectsClient)
}

func (c *client) BatchV1Beta1() batchv1beta1.Client {
	return batchv1beta1.NewClient(c.ObjectsClient)
}

func (c *client) BatchV2Alpha1() batchv2alpha1.Client {
	return batchv2alpha1.NewClient(c.ObjectsClient)
}

