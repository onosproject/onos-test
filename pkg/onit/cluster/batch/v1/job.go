package v1

import (
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

func newJob(object *clustermetav1.Object) *Job {
	return &Job{
		Object: object,
		Job: object.Object.(*batchv1.Job),
	}
}

type Job struct {
	*clustermetav1.Object
	Job *batchv1.Job
}
