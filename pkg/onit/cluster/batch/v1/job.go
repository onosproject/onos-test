package v1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var JobKind = clustermetav1.Kind{
	Group:   "batch",
	Version: "v1",
	Kind:    "Job",
}

var JobResource = clustermetav1.Resource{
	Kind: JobKind,
	Name: "Job",
	ObjectFactory: func() runtime.Object {
		return &batchv1.Job{}
	},
	ObjectsFactory: func() runtime.Object {
		return &batchv1.JobList{}
	},
}

type JobsClient interface {
	Get(name string) (*Job, error)
	List() ([]*Job, error)
}

// newJobsClient creates a new JobsClient
func newJobsClient(objects clustermetav1.ObjectsClient) JobsClient {
	return &jobsClient{
		ObjectsClient: objects,
	}
}

// jobsClient implements the JobsClient interface
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

// newJob creates a new Job resource
func newJob(object *clustermetav1.Object) *Job {
	return &Job{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*batchv1.Job).Spec,
	}
}

// Job provides functions for querying a job
type Job struct {
	*clustercorev1.PodSet
	Spec batchv1.JobSpec
}
