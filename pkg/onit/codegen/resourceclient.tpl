package {{ .Package.Name }}

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type {{ .Types.Interface }} interface {
	Get(name string) (*{{ .Resource.Types.Struct }}, error)
	List() ([]*{{ .Resource.Types.Struct }}, error)
}

func new{{ .Types.Interface }}(objects clustermetav1.ObjectsClient) {{ .Types.Interface }} {
	return &{{ .Types.Struct }}{
		ObjectsClient: objects,
	}
}

type {{ .Types.Struct }} struct {
	clustermetav1.ObjectsClient
}

func (c *{{ .Types.Struct }}) Get(name string) (*{{ .Resource.Types.Struct }}, error) {
	object, err := c.ObjectsClient.Get(name, {{ .Resource.Types.Resource }})
	if err != nil {
		return nil, err
	}
	return new{{ .Resource.Types.Struct }}(object), nil
}

func (c *{{ .Types.Struct }}) List() ([]*{{ .Resource.Types.Struct }}, error) {
	objects, err := c.ObjectsClient.List({{ .Resource.Types.Resource }})
	if err != nil {
		return nil, err
	}
	{{ .Resource.Names.Plural | toLowerCamel }} := make([]*{{ .Resource.Types.Struct }}, len(objects))
	for i, object := range objects {
		{{ .Resource.Names.Plural | toLowerCamel }}[i] = new{{ .Resource.Types.Struct }}(object)
	}
	return {{ .Resource.Names.Plural | toLowerCamel }}, nil
}
