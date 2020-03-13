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
    {{- range $name, $version := .Versions }}
    {{ $version.Package.Alias }} {{ $version.Package.Path | quote }}
    {{- end }}
    "github.com/onosproject/onos-test/pkg/onit/api/resource"
)

type {{ .Types.Interface }} interface {
    {{- $group := . }}
    {{- range $name, $version := .Versions }}
    {{ $group.Names.Proper }}{{ $version.Names.Proper }}() {{ $version.Package.Alias }}.{{ $version.Types.Interface }}
    {{- end }}
}

func New{{ .Types.Interface }}(resources resource.Client) {{ .Types.Interface }} {
	return &{{ .Types.Struct }}{
		Client: resources,
	}
}

type {{ .Types.Struct }} struct {
	resource.Client
}

{{ $group := . }}
{{- range $name, $version := .Versions }}
func (c *{{ $group.Types.Struct }}) {{ $group.Names.Proper }}{{ $version.Names.Proper }}() {{ $version.Package.Alias }}.{{ $version.Types.Interface }} {
	return {{ $version.Package.Alias }}.New{{ $version.Types.Interface }}(c.Client)
}
{{ end }}
