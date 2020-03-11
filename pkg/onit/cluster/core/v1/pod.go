package v1

import (
	"bytes"
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	executil "k8s.io/client-go/util/exec"
	"strings"
)

var PodKind = clustermetav1.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Pod",
}

var PodResource = clustermetav1.Resource{
	Kind: PodKind,
	Name: "Pod",
	ObjectFactory: func() runtime.Object {
		return &corev1.Pod{}
	},
	ObjectsFactory: func() runtime.Object {
		return &corev1.PodList{}
	},
}

type PodsClient interface {
	Get(name string) (*Pod, error)
	List() ([]*Pod, error)
}

// newPodsClient creates a new PodsClient
func newPodsClient(objects clustermetav1.ObjectsClient) PodsClient {
	return &podsClient{
		ObjectsClient: objects,
	}
}

// podsClient implements the PodsClient interface
type podsClient struct {
	clustermetav1.ObjectsClient
}

func (c *podsClient) Get(name string) (*Pod, error) {
	object, err := c.ObjectsClient.Get(name, PodResource)
	if err != nil {
		return nil, err
	}
	return newPod(object), nil
}

func (c *podsClient) List() ([]*Pod, error) {
	objects, err := c.ObjectsClient.List(PodResource)
	if err != nil {
		return nil, err
	}
	pods := make([]*Pod, len(objects))
	for i, object := range objects {
		pods[i] = newPod(object)
	}
	return pods, nil
}

// newPod creates a new Pod resource
func newPod(object *clustermetav1.Object) *Pod {
	return &Pod{
		Object: object,
		Spec:   object.Object.(*corev1.Pod).Spec,
	}
}

// Pod is a Kubernetes pod
type Pod struct {
	*clustermetav1.Object
	Spec corev1.PodSpec
}

// Execute executes the given command on the node
func (p *Pod) Execute(command ...string) ([]string, int, error) {
	container := p.Spec.Containers[0]
	fullCommand := append([]string{"/bin/bash", "-c"}, command...)
	req := p.Client().CoreV1().RESTClient().Post().
		Resource("pods").
		Name(p.Name).
		Namespace(p.Namespace).
		SubResource("exec").
		Param("container", container.Name)
	req.VersionedParams(&corev1.PodExecOptions{
		Container: container.Name,
		Command:   fullCommand,
		Stdout:    true,
		Stderr:    true,
		Stdin:     false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.Config(), "POST", req.URL())
	if err != nil {
		if execErr, ok := err.(executil.ExitError); ok && execErr.Exited() {
			return []string{}, execErr.ExitStatus(), nil
		}
		return nil, 0, err
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})

	if err != nil {
		if execErr, ok := err.(executil.ExitError); ok && execErr.Exited() {
			return []string{}, execErr.ExitStatus(), nil
		}
		return nil, 0, err
	}

	return strings.Split(strings.Trim(stdout.String(), "\n"), "\n"), 0, nil
}

// Kill kills the pod
func (p *Pod) Kill() error {
	return p.Client().CoreV1().Pods(p.Namespace).Delete(p.Name, &meta.DeleteOptions{})
}
