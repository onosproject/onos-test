package {{ .Package }}

import (
    {{- range $name, $group := .Groups }}
    "{{ $group.ImportPath }}"
    {{- end }}
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	"github.com/onosproject/onos-test/pkg/onit/cluster/networking"
)

{{- range $name, $group := .Groups }}
type {{ $group.Group }}Client {{ $group.Group }}.Client
{{- end }}

type Client interface {
    {{- range $name, $group := .Groups }}
    {{ $group.Group }}Client
    {{- end }}
}

func newClient(objects metav1.ObjectsClient) *client {
	return &client{
		ObjectsClient:    objects,
        {{- range $name, $group := .Groups }}
        {{ $group.Group }}Client: {{ $group.Group }}.NewClient(objects),
        {{- end }}
	}
}

type client struct {
	metav1.ObjectsClient
    {{- range $name, $group := .Groups }}
    {{ $group.Group }}Client
    {{- end }}
}
