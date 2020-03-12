package v1

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var ConfigMapKind = clustermetav1.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "ConfigMap",
}

var ConfigMapResource = clustermetav1.Resource{
	Kind: ConfigMapKind,
	Name: "ConfigMap",
	ObjectFactory: func() runtime.Object {
		return &corev1.ConfigMap{}
	},
	ObjectsFactory: func() runtime.Object {
		return &corev1.ConfigMapList{}
	},
}

func newConfigMap(object *clustermetav1.Object) *ConfigMap {
	return &ConfigMap{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*corev1.ConfigMap).Spec,
	}
}

type ConfigMap struct {
	*clustercorev1.PodSet
	Spec corev1.ConfigMapSpec
}
