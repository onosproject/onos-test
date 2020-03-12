package {{ .Package.Name }}

import (
    {{- range $name, $version := .Versions }}
    {{ $version.Package.Alias }} {{ $version.Package.Path | quote }}
    {{- end }}
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type {{ .Types.Interface }} interface {
    {{- range $name, $version := .Versions }}
    {{ .Names.Proper }}{{ $version.Names.Proper }}() {{ $version.Package.Alias }}.{{ $version.Types.Interface }}
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

{{- range $name, $version := .Versions }}
func (c *{{ .Types.Struct }}) {{ .Names.Proper }}{{ $version.Names.Proper }}() {{ $version.Package.Alias }}.{{ $version.Types.Interface }} {
	return {{ $version.Package.Alias }}.New{{ $version.Types.Interface }}(c.ObjectsClient)
}
{{- end }}
