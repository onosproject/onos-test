package v2alpha1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type CronJobsClient interface {
	Get(name string) (*CronJob, error)
	List() ([]*CronJob, error)
}

func newCronJobsClient(objects clustermetav1.ObjectsClient) CronJobsClient {
	return &cronJobsClient{
		ObjectsClient: objects,
	}
}

type cronJobsClient struct {
	clustermetav1.ObjectsClient
}

func (c *cronJobsClient) Get(name string) (*CronJob, error) {
	object, err := c.ObjectsClient.Get(name, CronJobResource)
	if err != nil {
		return nil, err
	}
	return newCronJob(object), nil
}

func (c *cronJobsClient) List() ([]*CronJob, error) {
	objects, err := c.ObjectsClient.List(CronJobResource)
	if err != nil {
		return nil, err
	}
	cronJobs := make([]*CronJob, len(objects))
	for i, object := range objects {
		cronJobs[i] = newCronJob(object)
	}
	return cronJobs, nil
}
