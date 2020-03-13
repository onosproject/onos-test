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

{{ $resource := .Resource }}
package {{ $resource.Package.Name }}

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	{{ $resource.Kind.Package.Alias }} {{ $resource.Kind.Package.Path | quote }}
    {{- range $sub := $resource.SubResources }}
    {{- if not (eq $sub.Client.Package.Path $resource.Package.Path) }}
    {{ $sub.Client.Package.Alias }} {{ $sub.Client.Package.Path | quote }}
    {{- end }}
    {{- end }}
	"k8s.io/apimachinery/pkg/runtime"
)

var {{ $resource.Types.Kind }} = clustermetav1.Kind{
	Group:   {{ $resource.Kind.Group | quote }},
	Version: {{ $resource.Kind.Version | quote }},
	Kind:    {{ $resource.Kind.Kind | quote }},
}

var {{ $resource.Types.Resource }} = clustermetav1.Resource{
	Kind: {{ $resource.Types.Kind }},
	Name: {{ $resource.Kind.Kind | quote }},
	ObjectFactory: func() runtime.Object {
		return &{{ $resource.Kind.Package.Alias }}.{{ $resource.Kind.Kind }}{}
	},
	ObjectsFactory: func() runtime.Object {
		return &{{ $resource.Kind.Package.Alias }}.{{ $resource.Kind.ListKind }}{}
	},
}

func New{{ $resource.Types.Struct }}(object *clustermetav1.Object) *{{ $resource.Types.Struct }} {
	return &{{ $resource.Types.Struct }}{
		Object: object,
		{{ $resource.Names.Singular }}: object.Object.(*{{ $resource.Kind.Package.Alias }}.{{ $resource.Kind.Kind }}),
        {{- range $sub := $resource.SubResources }}
        {{- if eq $sub.Resource.Package.Path $resource.Package.Path }}
        {{ $sub.Client.Types.Interface }}: New{{ $sub.Client.Types.Interface }}(object),
        {{- else }}
        {{ $sub.Client.Types.Interface }}: {{ $sub.Client.Package.Alias }}.New{{ $sub.Client.Types.Interface }}(object),
        {{- end }}
        {{- end }}
	}
}

type {{ $resource.Types.Struct }} struct {
	*clustermetav1.Object
	{{ $resource.Names.Singular }} *{{ $resource.Kind.Package.Alias }}.{{ $resource.Kind.Kind }}
    {{- range $sub := .Resource.SubResources }}
    {{- if eq $sub.Resource.Package.Path $resource.Package.Path }}
    {{ $sub.Client.Types.Interface }}
    {{- else }}
    {{ $sub.Client.Package.Alias }}.{{ $sub.Client.Types.Interface }}
    {{- end }}
    {{- end }}
}
