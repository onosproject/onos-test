package {{ .Package.Name }}

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	{{ .Kind.Package.Alias }} {{ .Kind.Package.Path | quote }}
	"k8s.io/apimachinery/pkg/runtime"
)

var {{ .Types.Kind }} = clustermetav1.Kind{
	Group:   {{ .Kind.Group | quote }},
	Version: {{ .Kind.Version | quote }},
	Kind:    {{ .Kind.Kind | quote }},
}

var {{ .Types.Resource }} = clustermetav1.Resource{
	Kind: {{ .Types.Kind }},
	Name: {{ .Kind.Kind | quote }},
	ObjectFactory: func() runtime.Object {
		return &{{ .Kind.Package.Alias }}.{{ .Kind.Kind }}{}
	},
	ObjectsFactory: func() runtime.Object {
		return &{{ .Kind.Package.Alias }}.{{ .Kind.ListKind }}{}
	},
}

func new{{ .Types.Struct }}(object *clustermetav1.Object) *{{ .Types.Struct }} {
	return &{{ .Types.Struct }}{
		Object: object,
		{{ .Names.Singular }}: object.Object.(*{{ .Kind.Package.Alias }}.{{ .Kind.Kind }}),
	}
}

type {{ .Types.Struct }} struct {
	*clustermetav1.Object
	{{ .Names.Singular }} *{{ .Kind.Package.Alias }}.{{ .Kind.Kind }}
}
