package v1

import (
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
	DaemonSets() DaemonSetsClient
	ReplicaSets() ReplicaSetsClient
	Deployments() DeploymentsClient
	StatefulSets() StatefulSetsClient
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}

func (c *client) DaemonSets() DaemonSetsClient {
	return newDaemonSetsClient(c.ObjectsClient)
}

func (c *client) ReplicaSets() ReplicaSetsClient {
	return newReplicaSetsClient(c.ObjectsClient)
}

func (c *client) Deployments() DeploymentsClient {
	return newDeploymentsClient(c.ObjectsClient)
}

func (c *client) StatefulSets() StatefulSetsClient {
	return newStatefulSetsClient(c.ObjectsClient)
}
