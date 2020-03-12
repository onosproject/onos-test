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

package {{ .Reader.Package.Name }}

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type {{ .Reader.Types.Interface }} interface {
	Get(name string) (*{{ .Resource.Types.Struct }}, error)
	List() ([]*{{ .Resource.Types.Struct }}, error)
}

func New{{ .Reader.Types.Interface }}(objects clustermetav1.ObjectsClient) {{ .Reader.Types.Interface }} {
	return &{{ .Reader.Types.Struct }}{
		ObjectsClient: objects,
	}
}

type {{ .Reader.Types.Struct }} struct {
	clustermetav1.ObjectsClient
}

func (c *{{ .Reader.Types.Struct }}) Get(name string) (*{{ .Resource.Types.Struct }}, error) {
	object, err := c.ObjectsClient.Get(name, {{ .Resource.Types.Resource }})
	if err != nil {
		return nil, err
	}
	return New{{ .Resource.Types.Struct }}(object), nil
}

func (c *{{ .Reader.Types.Struct }}) List() ([]*{{ .Resource.Types.Struct }}, error) {
	objects, err := c.ObjectsClient.List({{ .Resource.Types.Resource }})
	if err != nil {
		return nil, err
	}
	{{ .Resource.Names.Plural | toLowerCamel }} := make([]*{{ .Resource.Types.Struct }}, len(objects))
	for i, object := range objects {
		{{ .Resource.Names.Plural | toLowerCamel }}[i] = New{{ .Resource.Types.Struct }}(object)
	}
	return {{ .Resource.Names.Plural | toLowerCamel }}, nil
}
