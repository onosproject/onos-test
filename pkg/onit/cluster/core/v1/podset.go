package v1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

func NewPodSet(object *clustermetav1.Object) *PodSet {
	return &PodSet{
		Object: object,
	}
}

// PodSet provides functions for querying a collection of pods
type PodSet struct {
	*clustermetav1.Object
}

func (s *PodSet) GetPod(name string) (*Pod, error) {
	object, err := s.ObjectsClient.Get(name, PodResource)
	if err != nil {
		return nil, err
	}
	return newPod(object), nil
}

func (s *PodSet) GetPods() ([]*Pod, error) {
	objects, err := s.ObjectsClient.List(PodResource)
	if err != nil {
		return nil, err
	}
	pods := make([]*Pod, len(objects))
	for i, object := range objects {
		pods[i] = newPod(object)
	}
	return pods, nil
}
