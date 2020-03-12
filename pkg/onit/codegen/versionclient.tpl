package {{ .Package.Name }}

import (
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type {{ .Types.Interface }} interface {
    {{- range $name, $resource := .Resources }}
    {{ $resource.PluralKind }}() {{ $resource.PluralKind }}Client
    {{- end }}
}

func New{{ .Types.Interface }}(objects metav1.ObjectsClient) {{ .Types.Interface }} {
	return &{{ .Types.Struct }}{
		ObjectsClient: objects,
	}
}

type {{ .Types.Struct }} struct {
	metav1.ObjectsClient
}

{{- range $name, $resource := .Resources }}
func (c *{{ .Types.Struct }}) {{ $resource.Names.Plural }}() {{ $resource.Types.Interface }} {
	return new{{ $resource.Types.Interface }}(c.ObjectsClient)
}
{{- end }}
