package {{ .Package }}

import (
    {{- range $name, $version := .Versions }}
    "{{ $version.ImportPath }}"
    {{- end }}
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
    {{- range $name, $version := .Versions }}
    {{ $version.Group | toCamelCase }}{{ $version.Version | toCamelCase }}() {{ $version.Version }}.Client
    {{- end }}
}

func NewClient(objects metav1.ObjectsClient) Client {
	return &client{
		ObjectsClient: objects,
	}
}

type client struct {
	metav1.ObjectsClient
}

{{- range $name, $version := .Versions }}
func (c *client) {{ $version.Group | toCamelCase }}{{ $version.Version | toCamelCase }}() {{ $version.Version }}.Client {
	return {{ $version.Version }}.NewClient(c.ObjectsClient)
}
{{- end }}
