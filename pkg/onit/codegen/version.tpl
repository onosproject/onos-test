package {{ .Package }}

import (
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type Client interface {
    {{- range $name, $resource := .Resources }}
    {{ $resource.PluralKind }}() {{ $resource.PluralKind }}Client
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

{{- range $name, $resource := .Resources }}
func (c *client) {{ $resource.PluralKind }}() {{ $resource.PluralKind }}Client {
	return new{{ $resource.PluralKind }}Client(c.ObjectsClient)
}
{{- end }}
