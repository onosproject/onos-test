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
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type {{ .Types.Interface }} interface {
    {{- range $name, $resource := .Resources }}
    {{ $resource.Client.Types.Interface }}
    {{- end }}
}

func New{{ .Types.Interface }}(objects metav1.ObjectsClient) {{ .Types.Interface }} {
	return &{{ .Types.Struct }}{
		ObjectsClient: objects,
		{{- range $name, $resource := .Resources }}
        {{ $resource.Client.Types.Interface }}: New{{ $resource.Client.Types.Interface }}(objects),
        {{ end }}
	}
}

type {{ .Types.Struct }} struct {
	metav1.ObjectsClient
    {{- range $name, $resource := .Resources }}
    {{ $resource.Client.Types.Interface }}
    {{- end }}
}
