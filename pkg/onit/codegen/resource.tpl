// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

{{- $resource := .Resource }}
{{- $field := .Resource.Names.Singular }}
{{- $name := (.Resource.Names.Singular | toLowerCamel) }}
{{- $kind := (printf "%s.%s" .Resource.Kind.Package.Alias .Resource.Kind.Kind) }}

package {{ $resource.Package.Name }}

import (
    "github.com/onosproject/onos-test/pkg/onit/api/resource"
	{{ .Resource.Kind.Package.Alias }} {{ .Resource.Kind.Package.Path | quote }}
    {{- range $ref := $resource.References }}
    {{- if not (eq $ref.Reference.Package.Path $resource.Package.Path) }}
    {{ $ref.Reference.Package.Alias }} {{ $ref.Reference.Package.Path | quote }}
    {{- end }}
    {{- end }}
)

var {{ $resource.Types.Kind }} = resource.Kind{
	Group:   {{ $resource.Kind.Group | quote }},
	Version: {{ $resource.Kind.Version | quote }},
	Kind:    {{ $resource.Kind.Kind | quote }},
}

var {{ $resource.Types.Resource }} = resource.Type{
	Kind: {{ $resource.Types.Kind }},
	Name: {{ $resource.Names.Plural | lower | quote }},
}

func New{{ $resource.Types.Struct }}({{ $name }} *{{ $kind }}, client resource.Client) *{{ $resource.Types.Struct }} {
	return &{{ $resource.Types.Struct }}{
		Resource: resource.NewResource({{ $name }}.ObjectMeta, {{ .Resource.Types.Kind }}, client),
		{{ $field }}: {{ $name }},
        {{- range $ref := $resource.References }}
        {{- if eq $ref.Resource.Package.Path $resource.Package.Path }}
        {{ $ref.Reference.Types.Interface }}: New{{ $ref.Reference.Types.Interface }}(client, resource.NewUIDFilter({{ $name }}.UID)),
        {{- else }}
        {{ $ref.Reference.Types.Interface }}: {{ $ref.Reference.Package.Alias }}.New{{ $ref.Reference.Types.Interface }}(client, resource.NewUIDFilter({{ $name }}.UID)),
        {{- end }}
        {{- end }}
	}
}

type {{ $resource.Types.Struct }} struct {
	*resource.Resource
	{{ $field }} *{{ $kind }}
    {{- range $ref := .Resource.References }}
    {{- if eq $ref.Resource.Package.Path $resource.Package.Path }}
    {{ $ref.Reference.Types.Interface }}
    {{- else }}
    {{ $ref.Reference.Package.Alias }}.{{ $ref.Reference.Types.Interface }}
    {{- end }}
    {{- end }}
}
