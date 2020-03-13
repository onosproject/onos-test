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
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

{{- range $name, $group := .Groups }}
type {{ $group.Group }}Client {{ $group.Package.Alias }}.{{ $group.Types.Interface }}
{{- end }}

type {{ .Types.Interface }} interface {
    {{- range $name, $group := .Groups }}
    {{ $group.Group }}Client
    {{- end }}
}

func New{{ .Types.Interface }}(objects metav1.ObjectsClient) {{ .Types.Interface }} {
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
