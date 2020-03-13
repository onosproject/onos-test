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

package {{ .Reference.Package.Name }}

import (
    "github.com/onosproject/onos-test/pkg/onit/api/resource"
    {{- $resource := .Resource }}
    {{- range $ref := $resource.References }}
    {{- if not (eq $ref.Reference.Package.Path $resource.Package.Path) }}
    {{ $ref.Reference.Package.Alias }} {{ $ref.Reference.Package.Path | quote }}
    {{- end }}
    {{- end }}
    {{- if .Resource.References }}
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	{{- end }}
)

type {{ .Reference.Types.Interface }} interface {
    {{ .Resource.Names.Plural }}() {{ .Reader.Types.Interface }}
    {{- range $ref := .Resource.References }}
    {{- if eq $ref.Resource.Package.Path $resource.Package.Path }}
    {{ $ref.Reference.Types.Interface }}
    {{- else }}
    {{ $ref.Reference.Package.Alias }}.{{ $ref.Reference.Types.Interface }}
    {{- end }}
    {{- end }}
}

func New{{ .Reference.Types.Interface }}(resources resource.Client, filter resource.Filter) {{ .Reference.Types.Interface }} {
    {{- if .Resource.References }}
    var ownerFilter resource.Filter = func(kind metav1.GroupVersionKind, meta metav1.ObjectMeta) (bool, error) {
        {{- $name :=  (.Resource.Names.Plural | toLowerCamel) }}
        list, err := New{{ .Reader.Types.Interface }}(resources, filter).List()
        if err != nil {
            return false, err
        }
        for _, {{ $name }} := range list {
            if {{ $name }}.{{ .Resource.Names.Singular }}.ObjectMeta.UID == meta.UID {
                return true, nil
            }
        }
        return false, nil
    }
    {{- end }}
	return &{{ .Reference.Types.Struct }}{
		Client: resources,
		filter: filter,
        {{- range $ref := $resource.References }}
        {{- if eq $ref.Resource.Package.Path $resource.Package.Path }}
        {{ $ref.Reference.Types.Interface }}: New{{ $ref.Reference.Types.Interface }}(resources, ownerFilter),
        {{- else }}
        {{ $ref.Reference.Types.Interface }}: {{ $ref.Reference.Package.Alias }}.New{{ $ref.Reference.Types.Interface }}(resources, ownerFilter),
        {{- end }}
        {{- end }}
	}
}

type {{ .Reference.Types.Struct }} struct {
	resource.Client
	filter resource.Filter
    {{- range $ref := .Resource.References }}
    {{- if eq $ref.Resource.Package.Path $resource.Package.Path }}
    {{ $ref.Reference.Types.Interface }}
    {{- else }}
    {{ $ref.Reference.Package.Alias }}.{{ $ref.Reference.Types.Interface }}
    {{- end }}
    {{- end }}
}

func (c *{{ .Reference.Types.Struct }}) {{ .Resource.Names.Plural }}() {{ .Reader.Types.Interface }} {
    return New{{ .Reader.Types.Interface }}(c.Client, c.filter)
}
