package {{ .Package }}

import (
	clustercorev1 "github.com/onosproject/onos-test/pkg/onit/cluster/core/v1"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	{{ .Group }}{{ .Version }} "k8s.io/api/{{ .Group }}/{{ .Version }}"
	"k8s.io/apimachinery/pkg/runtime"
)

var {{ .Kind }}Kind = clustermetav1.Kind{
	Group:   "{{ .Group }}",
	Version: "{{ .Version }}",
	Kind:    "{{ .Kind }}",
}

var {{ .Kind }}Resource = clustermetav1.Resource{
	Kind: {{ .Kind }}Kind,
	Name: "{{ .Kind }}",
	ObjectFactory: func() runtime.Object {
		return &{{ .Group }}{{ .Version }}.{{ .Kind }}{}
	},
	ObjectsFactory: func() runtime.Object {
		return &{{ .Group }}{{ .Version }}.{{ .ListKind }}{}
	},
}

type {{ .PluralKind }}Client interface {
	Get(name string) (*{{ .Kind }}, error)
	List() ([]*{{ .Kind }}, error)
}

func new{{ .PluralKind }}Client(objects clustermetav1.ObjectsClient) {{ .PluralKind }}Client {
	return &{{ .Plural }}Client{
		ObjectsClient: objects,
	}
}

type {{ .Plural }}Client struct {
	clustermetav1.ObjectsClient
}

func (c *{{ .Plural }}Client) Get(name string) (*{{ .Kind }}, error) {
	object, err := c.ObjectsClient.Get(name, {{ .Kind }}Resource)
	if err != nil {
		return nil, err
	}
	return new{{ .Kind }}(object), nil
}

func (c *{{ .Plural }}Client) List() ([]*{{ .Kind }}, error) {
	objects, err := c.ObjectsClient.List({{ .Kind }}Resource)
	if err != nil {
		return nil, err
	}
	{{ .Plural }} := make([]*{{ .Kind }}, len(objects))
	for i, object := range objects {
		{{ .Plural }}[i] = new{{ .Kind }}(object)
	}
	return {{ .Plural }}, nil
}

func new{{ .Kind }}(object *clustermetav1.Object) *{{ .Kind }} {
	return &{{ .Kind }}{
		PodSet: clustercorev1.NewPodSet(object),
		Spec:   object.Object.(*{{ .Group }}{{ .Version }}.{{ .Kind }}).Spec,
	}
}

type {{ .Kind }} struct {
	*clustercorev1.PodSet
	Spec {{ .Group }}{{ .Version }}.{{ .Kind }}Spec
}
