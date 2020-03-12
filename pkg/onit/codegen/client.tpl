package {{ .Package.Name }}

import (
    {{- range $name, $group := .Groups }}
    {{ $group.Package.Alias }} {{ $group.Package.Path | quote }}
    {{- end }}
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	"github.com/onosproject/onos-test/pkg/onit/cluster/networking"
)

{{- range $name, $group := .Groups }}
type {{ $group.Group }}Client {{ $group.Package.Alias }}.{{ $group.Types.Interface }}
{{- end }}

type {{ .Types.Interface }} interface {
    {{- range $name, $group := .Groups }}
    {{ $group.Group }}Client
    {{- end }}
}

func new{{ .Types.Interface }}(objects metav1.ObjectsClient) *{{ .Types.Struct }} {
	return &{{ .Types.Struct }}{
		ObjectsClient:    objects,
        {{- range $name, $group := .Groups }}
        {{ $group.Group }}Client: {{ $group.Package.Alias }}.New{{ $group.Types.Interface }}(objects),
        {{- end }}
	}
}

type {{ .Types.Struct }} struct {
	metav1.ObjectsClient
    {{- range $name, $group := .Groups }}
    {{ $group.Group }}Client
    {{- end }}
}
