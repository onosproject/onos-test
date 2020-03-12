package v1beta1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var CronJobKind = clustermetav1.Kind{
	Group:   "batch",
	Version: "v1beta1",
	Kind:    "CronJob",
}

var CronJobResource = clustermetav1.Resource{
	Kind: CronJobKind,
	Name: "CronJob",
	ObjectFactory: func() runtime.Object {
		return &batchv1beta1.CronJob{}
	},
	ObjectsFactory: func() runtime.Object {
		return &batchv1beta1.CronJobList{}
	},
}

func newCronJob(object *clustermetav1.Object) *CronJob {
	return &CronJob{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*batchv1beta1.CronJob).Spec,
	}
}

type CronJob struct {
	*clustercorev1.PodSet
	Spec batchv1beta1.CronJobSpec
}
