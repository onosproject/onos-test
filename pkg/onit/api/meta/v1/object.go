package v1

import (
	"errors"
	"fmt"
	"github.com/onosproject/onos-test/pkg/kube"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"time"
)

const apiTimeout = 30 * time.Second

// Resource is a resource
type Resource struct {
	Kind           Kind
	Name           string
	ObjectFactory  ObjectFactory
	ObjectsFactory ObjectsFactory
}

// Kind is a resource kind
type Kind struct {
	Group   string
	Version string
	Kind    string
}

func (k Kind) APIVersion() string {
	if k.Group == "core" {
		return k.Version
	}
	return fmt.Sprintf("%s/%s", k.Group, k.Version)
}

func getTypeMeta(kind Kind) metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: kind.APIVersion(),
		Kind:       kind.Kind,
	}
}

// ObjectFilter is an object filter
type ObjectFilter func(Object) (bool, error)

// NewCompositeFilter returns a composite object filter
func NewCompositeFilter(filters ...ObjectFilter) ObjectFilter {
	return func(object Object) (bool, error) {
		for _, filter := range filters {
			ok, err := filter(object)
			if !ok || err != nil {
				return ok, err
			}
		}
		return true, nil
	}
}

// ObjectFactory is an object factory
type ObjectFactory func() runtime.Object

// ObjectsFactory is an object list factory
type ObjectsFactory func() runtime.Object

func NewObjectsClient(api kube.API, filter ObjectFilter) ObjectsClient {
	return &objectsClient{
		API:     api,
		isValid: filter,
	}
}

// ObjectsClient is an interface for querying objects in the cluster
type ObjectsClient interface {
	kube.API
	Get(name string, resource Resource) (*Object, error)
	List(resource Resource) ([]*Object, error)
}

// objectsClient is an implementation of the ObjectsClient interface
type objectsClient struct {
	kube.API
	isValid ObjectFilter
}

func (c *objectsClient) newObjects(runtimeObjects runtime.Object) ([]*Object, error) {
	elem := reflect.ValueOf(runtimeObjects).Elem()
	elemType := elem.Type()
	field, ok := elemType.FieldByName("Items")
	if !ok {
		return nil, errors.New("unable to locate object metadata items")
	}
	items := elem.FieldByIndex(field.Index)

	objects := make([]*Object, items.Len())
	for i := 0; i < items.Len(); i++ {
		object, err := c.newObject(items.Index(i).Interface().(runtime.Object))
		if err != nil {
			return nil, err
		}
		objects[i] = object
	}
	return objects, nil
}

func getObjectMeta(object runtime.Object) (metav1.ObjectMeta, error) {
	elem := reflect.ValueOf(object).Elem()
	elemType := elem.Type()
	field, ok := elemType.FieldByName("ObjectMeta")
	if !ok {
		return metav1.ObjectMeta{}, errors.New("unable to locate object metadata")
	}
	return elem.FieldByIndex(field.Index).Interface().(metav1.ObjectMeta), nil
}

func (c *objectsClient) newObject(runtimeObject runtime.Object) (*Object, error) {
	meta, err := getObjectMeta(runtimeObject)
	if err != nil {
		return nil, err
	}
	kind := Kind{
		Group:   runtimeObject.GetObjectKind().GroupVersionKind().Group,
		Version: runtimeObject.GetObjectKind().GroupVersionKind().Version,
		Kind:    runtimeObject.GetObjectKind().GroupVersionKind().Kind,
	}
	filter := NewCompositeFilter(c.isValid, func(object Object) (bool, error) {
		filterMeta, err := getObjectMeta(object.Object)
		if err != nil {
			return false, err
		}
		for _, owner := range filterMeta.OwnerReferences {
			if owner.UID == meta.UID {
				return true, nil
			}
		}
		return false, nil
	})
	client := NewObjectsClient(c, filter)
	return NewObject(runtimeObject, meta, kind, client), nil
}

func (c *objectsClient) Get(name string, resource Resource) (*Object, error) {
	runtimeObject := resource.ObjectFactory()
	opts := &metav1.GetOptions{
		TypeMeta: getTypeMeta(resource.Kind),
	}
	err := c.Client().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Name(name).
		VersionedParams(opts, metav1.ParameterCodec).
		Timeout(apiTimeout).
		Do().
		Into(runtimeObject)
	if err != nil {
		return nil, err
	}

	object, err := c.newObject(runtimeObject)
	if err != nil {
		return nil, k8serrors.NewInternalError(err)
	}

	ok, err := c.isValid(*object)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, k8serrors.NewNotFound(schema.GroupResource{
			Group:    resource.Kind.Group,
			Resource: resource.Name,
		}, name)
	}
	return object, nil
}

func (c *objectsClient) List(resource Resource) ([]*Object, error) {
	runtimeObjects := resource.ObjectsFactory()
	opts := &metav1.ListOptions{
		TypeMeta: getTypeMeta(resource.Kind),
	}
	err := c.Client().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		VersionedParams(opts, metav1.ParameterCodec).
		Timeout(apiTimeout).
		Do().
		Into(runtimeObjects)
	if err != nil {
		return nil, err
	}

	objects, err := c.newObjects(runtimeObjects)
	if err != nil {
		return nil, k8serrors.NewInternalError(err)
	}

	filtered := make([]*Object, 0)
	for _, object := range objects {
		ok, err := c.isValid(*object)
		if err != nil {
			return nil, err
		} else if ok {
			filtered = append(filtered, object)
		}
	}
	return filtered, nil
}

// NewObject creates a new object
func NewObject(object runtime.Object, meta metav1.ObjectMeta, kind Kind, client ObjectsClient) *Object {
	return &Object{
		ObjectsClient: client,
		Object:        object,
		Kind:          kind,
		Namespace:     meta.Namespace,
		Name:          meta.Name,
		UID:           meta.UID,
		Labels:        meta.Labels,
		Annotations:   meta.Annotations,
	}
}

// Object is a Kubernetes object
type Object struct {
	ObjectsClient
	Object      runtime.Object
	Kind        Kind
	Namespace   string
	Name        string
	UID         types.UID
	Labels      map[string]string
	Annotations map[string]string
}
