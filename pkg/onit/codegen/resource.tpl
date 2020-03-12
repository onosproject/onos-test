package {{ .Package.Name }}

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	{{ .Kind.Package.Alias }} {{ .Kind.Package.Path | quote }}
	"k8s.io/apimachinery/pkg/runtime"
)

var {{ .Types.Kind }} = clustermetav1.Kind{
	Group:   "{{ .Group }}",
	Version: "{{ .Version }}",
	Kind:    "{{ .Kind }}",
}

var {{ .Types.Resource }} = clustermetav1.Resource{
	Kind: {{ .Types.Kind }},
	Name: "{{ .Kind.Kind }}",
	ObjectFactory: func() runtime.Object {
		return &{{ .Kind.Package.Alias }}.{{ .Kind.Kind }}{}
	},
	ObjectsFactory: func() runtime.Object {
		return &{{ .Kind.Package.Alias }}.{{ .Kind.ListKind }}{}
	},
}

func new{{ .Types.Struct }}(object *clustermetav1.Object) *{{ .Types.Struct }} {
	return &{{ .Types.Struct }}{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*{{ .Kind.Package.Alias }}.{{ .Kind.Kind }}).Spec,
	}
}

type {{ .Types.Struct }} struct {
	*clustercorev1.PodSet
	Spec {{ .Kind.Package.Alias }}.{{ .Kind.Kind }}Spec
}
