package apps

import (
    appsv1 "github.com/onosproject/onos-test/pkg/onit/cluster/apps/v1"
    appsv1beta1 "github.com/onosproject/onos-test/pkg/onit/cluster/apps/v1beta1"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
    AppsV1() appsv1.Client
    AppsV1Beta1() appsv1beta1.Client
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}


func (c *client) AppsV1() appsv1.Client {
	return appsv1.NewClient(c.ObjectsClient)
}

func (c *client) AppsV1Beta1() appsv1beta1.Client {
	return appsv1beta1.NewClient(c.ObjectsClient)
}

