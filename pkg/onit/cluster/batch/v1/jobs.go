package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type JobsClient interface {
	Get(name string) (*Job, error)
	List() ([]*Job, error)
}

func newJobsClient(objects clustermetav1.ObjectsClient) JobsClient {
	return &jobsClient{
		ObjectsClient: objects,
	}
}

type jobsClient struct {
	clustermetav1.ObjectsClient
}

func (c *jobsClient) Get(name string) (*Job, error) {
	object, err := c.ObjectsClient.Get(name, JobResource)
	if err != nil {
		return nil, err
	}
	return newJob(object), nil
}

func (c *jobsClient) List() ([]*Job, error) {
	objects, err := c.ObjectsClient.List(JobResource)
	if err != nil {
		return nil, err
	}
	jobs := make([]*Job, len(objects))
	for i, object := range objects {
		jobs[i] = newJob(object)
	}
	return jobs, nil
}
