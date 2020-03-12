package v1beta1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var DeploymentKind = clustermetav1.Kind{
	Group:   "apps",
	Version: "v1beta1",
	Kind:    "Deployment",
}

var DeploymentResource = clustermetav1.Resource{
	Kind: DeploymentKind,
	Name: "Deployment",
	ObjectFactory: func() runtime.Object {
		return &appsv1beta1.Deployment{}
	},
	ObjectsFactory: func() runtime.Object {
		return &appsv1beta1.DeploymentList{}
	},
}

func newDeployment(object *clustermetav1.Object) *Deployment {
	return &Deployment{
		Object: object,
		Deployment: object.Object.(*appsv1beta1.Deployment),
	}
}

type Deployment struct {
	*clustermetav1.Object
	Deployment *appsv1beta1.Deployment
}
