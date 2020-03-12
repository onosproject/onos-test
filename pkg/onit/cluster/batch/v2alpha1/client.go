package v2alpha1

import (
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
    CronJobs() CronJobsClient
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}


func (c *client) CronJobs() CronJobsClient {
	return newCronJobsClient(c.ObjectsClient)
}

