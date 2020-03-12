package v2alpha1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

var CronJobKind = clustermetav1.Kind{
	Group:   "batch",
	Version: "v2alpha1",
	Kind:    "CronJob",
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

func newCronJob(object *clustermetav1.Object) *CronJob {
	return &CronJob{
		Object: object,
		CronJob: object.Object.(*batchv2alpha1.CronJob),
	}
}

type CronJob struct {
	*clustermetav1.Object
	CronJob *batchv2alpha1.CronJob
}
