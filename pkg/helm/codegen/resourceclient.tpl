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

package {{ .Client.Package.Name }}

import (
    "github.com/onosproject/onos-test/pkg/helm/api/resource"
)

type {{ .Client.Types.Interface }} interface {
    {{ .Resource.Names.Plural }}() {{ .Reader.Types.Interface }}
}

func New{{ .Client.Types.Interface }}(resources resource.Client, filter resource.Filter) {{ .Client.Types.Interface }} {
	return &{{ .Client.Types.Struct }}{
		Client: resources,
		filter: filter,
	}
}

type {{ .Client.Types.Struct }} struct {
	resource.Client
	filter resource.Filter
}

func (c *{{ .Client.Types.Struct }}) {{ .Resource.Names.Plural }}() {{ .Reader.Types.Interface }} {
    return New{{ .Reader.Types.Interface }}(c.Client, c.filter)
}
