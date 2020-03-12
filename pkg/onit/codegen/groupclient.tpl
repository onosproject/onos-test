package {{ .Package.Name }}

import (
    {{- range $name, $version := .Versions }}
    {{ $version.Package.Alias }} {{ $version.Package.Path | quote }}
    {{- end }}
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type {{ .Types.Interface }} interface {
    {{- $group := . }}
    {{- range $name, $version := .Versions }}
    {{ $group.Names.Proper }}{{ $version.Names.Proper }}() {{ $version.Package.Alias }}.{{ $version.Types.Interface }}
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

{{ $group := . }}
{{- range $name, $version := .Versions }}
func (c *{{ $group.Types.Struct }}) {{ $group.Names.Proper }}{{ $version.Names.Proper }}() {{ $version.Package.Alias }}.{{ $version.Types.Interface }} {
	return {{ $version.Package.Alias }}.New{{ $version.Types.Interface }}(c.ObjectsClient)
}
{{ end }}
