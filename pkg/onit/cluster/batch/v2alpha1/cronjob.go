package v2alpha1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

var CronJobKind = clustermetav1.Kind{
	Group:   "batch",
	Version: "v2alpha1",
	Kind:    "Job",
}

var CronJobResource = clustermetav1.Resource{
	Kind: CronJobKind,
	Name: "CronJob",
	ObjectFactory: func() runtime.Object {
		return &batchv2alpha1.CronJob{}
	},
	ObjectsFactory: func() runtime.Object {
		return &batchv2alpha1.CronJobList{}
	},
}

type CronJobsClient interface {
	Get(name string) (*CronJob, error)
	List() ([]*CronJob, error)
}

// newCronJobsClient creates a new CronJobsClient
func newCronJobsClient(objects clustermetav1.ObjectsClient) CronJobsClient {
	return &cronJobsClient{
		ObjectsClient: objects,
	}
}

// cronJobsClient implements the CronJobsClient interface
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

// newCronJob creates a new CronJob resource
func newCronJob(object *clustermetav1.Object) *CronJob {
	return &CronJob{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*batchv2alpha1.CronJob).Spec,
	}
}

// CronJob provides functions for querying a cron job
type CronJob struct {
	*clustercorev1.PodSet
	Spec batchv2alpha1.CronJobSpec
}
