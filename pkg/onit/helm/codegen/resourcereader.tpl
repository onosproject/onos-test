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
    "github.com/onosproject/onos-test/pkg/onit/helm/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	{{ .Resource.Kind.Package.Alias }} {{ .Resource.Kind.Package.Path | quote }}
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

type {{ .Reader.Types.Interface }} interface {
	Get(name string) (*{{ .Resource.Types.Struct }}, error)
	List() ([]*{{ .Resource.Types.Struct }}, error)
}

func New{{ .Reader.Types.Interface }}(client resource.Client, filter resource.Filter) {{ .Reader.Types.Interface }} {
	return &{{ .Reader.Types.Struct }}{
		Client: client,
		filter: filter,
	}
}

type {{ .Reader.Types.Struct }} struct {
	resource.Client
	filter resource.Filter
}

{{- $singular := (.Resource.Names.Singular | toLowerCamel) }}
{{- $kind := (printf "%s.%s" .Resource.Kind.Package.Alias .Resource.Kind.Kind) }}
{{- $listKind := (printf "%s.%s" .Resource.Kind.Package.Alias .Resource.Kind.ListKind) }}

func (c *{{ .Reader.Types.Struct }}) Get(name string) (*{{ .Resource.Types.Struct }}, error) {
    {{ $singular }} := &{{ $kind }}{}
	err := c.Clientset().
        {{ .Group.Names.Proper }}{{ .Version.Names.Proper }}().
        RESTClient().
	    Get().
		Namespace(c.Namespace()).
		Resource({{ .Resource.Types.Resource }}.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into({{ $singular }})
	if err != nil {
		return nil, err
	} else {
        ok, err := c.filter(metav1.GroupVersionKind{
            Group:   {{ .Resource.Types.Kind }}.Group,
            Version: {{ .Resource.Types.Kind }}.Version,
            Kind:    {{ .Resource.Types.Kind }}.Kind,
        }, {{ $singular }}.ObjectMeta)
        if err != nil {
            return nil, err
        } else if !ok {
            return nil, errors.NewNotFound(schema.GroupResource{
                Group:    {{ .Resource.Types.Kind }}.Group,
                Resource: {{ .Resource.Types.Resource }}.Name,
            }, name)
        }
    }
	return New{{ .Resource.Types.Struct }}({{ $singular }}, c.Client), nil
}

func (c *{{ .Reader.Types.Struct }}) List() ([]*{{ .Resource.Types.Struct }}, error) {
    list := &{{ $listKind }}{}
	err := c.Clientset().
        {{ .Group.Names.Proper }}{{ .Version.Names.Proper }}().
        RESTClient().
	    Get().
		Namespace(c.Namespace()).
		Resource({{ .Resource.Types.Resource }}.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*{{ .Resource.Types.Struct }}, 0, len(list.Items))
	for _, {{ $singular }} := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   {{ .Resource.Types.Kind }}.Group,
			Version: {{ .Resource.Types.Kind }}.Version,
			Kind:    {{ .Resource.Types.Kind }}.Kind,
		}, {{ $singular }}.ObjectMeta)
        if err != nil {
            return nil, err
        } else if ok {
    	    results = append(results, New{{ .Resource.Types.Struct }}(&{{ $singular }}, c.Client))
        }
	}
	return results, nil
}
