package v1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var DeploymentKind = clustermetav1.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "Deployment",
}

var DeploymentResource = clustermetav1.Resource{
	Kind: DeploymentKind,
	Name: "Deployment",
	ObjectFactory: func() runtime.Object {
		return &appsv1.Deployment{}
	},
	ObjectsFactory: func() runtime.Object {
		return &appsv1.DeploymentList{}
	},
}

func newDeployment(object *clustermetav1.Object) *Deployment {
	return &Deployment{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*appsv1.Deployment).Spec,
	}
}

type Deployment struct {
	*clustercorev1.PodSet
	Spec appsv1.DeploymentSpec
}
