// Code generated by onit-generate. DO NOT EDIT.

package v1beta1

import (
	appsv1 "github.com/onosproject/onos-test/pkg/helm/api/apps/v1"
	corev1 "github.com/onosproject/onos-test/pkg/helm/api/core/v1"
	"github.com/onosproject/onos-test/pkg/helm/api/resource"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var StatefulSetKind = resource.Kind{
	Group:   "apps",
	Version: "v1beta1",
	Kind:    "StatefulSet",
}

var StatefulSetResource = resource.Type{
	Kind: StatefulSetKind,
	Name: "statefulsets",
}

func NewStatefulSet(statefulSet *appsv1beta1.StatefulSet, client resource.Client) *StatefulSet {
	return &StatefulSet{
		Resource:             resource.NewResource(statefulSet.ObjectMeta, StatefulSetKind, client),
		Object:               statefulSet,
		ReplicaSetsReference: appsv1.NewReplicaSetsReference(client, resource.NewUIDFilter(statefulSet.UID)),
		PodsReference:        corev1.NewPodsReference(client, resource.NewUIDFilter(statefulSet.UID)),
	}
}

type StatefulSet struct {
	*resource.Resource
	Object *appsv1beta1.StatefulSet
	appsv1.ReplicaSetsReference
	corev1.PodsReference
}

func (r *StatefulSet) Delete() error {
	return r.Clientset().
		AppsV1beta1().
		RESTClient().
		Delete().
		Namespace(r.Namespace).
		Resource(StatefulSetResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}
