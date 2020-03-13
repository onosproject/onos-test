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

package {{ .Package.Name }}

import (
    {{- range $name, $group := .Groups }}
    {{ $group.Package.Alias }} {{ $group.Package.Path | quote }}
    {{- end }}
    "github.com/onosproject/onos-test/pkg/onit/helm/api/resource"
)

type {{ .Types.Interface }} interface {
    resource.Client
    {{- range $name, $group := .Groups }}
    {{ $group.Names.Proper }}() {{ $group.Package.Alias }}.{{ $group.Types.Interface }}
    {{- end }}
}

func New{{ .Types.Interface }}(resources resource.Client, filter resource.Filter) {{ .Types.Interface }} {
	return &{{ .Types.Struct }}{
		Client: resources,
		filter: filter,
	}
}

type {{ .Types.Struct }} struct {
	resource.Client
	filter resource.Filter
}

{{- range $name, $group := .Groups }}
func (c *{{ .Types.Struct }}) {{ $group.Names.Proper }}() {{ $group.Package.Alias }}.{{ $group.Types.Interface }} {
    return {{ $group.Package.Alias }}.New{{ $group.Types.Interface }}(c.Client, c.filter)
}
{{ end }}
