package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

func (d *Deployment) Wait(timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		deployment, err := d.Clientset().AppsV1().Deployments(d.Namespace).Get(d.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if deployment.Spec.Paused {
			return false, nil
		}
		if deployment.Spec.Strategy.RollingUpdate != nil && deployment.Spec.Strategy.RollingUpdate.MaxUnavailable != nil {
			return deployment.Status.UnavailableReplicas <= deployment.Spec.Strategy.RollingUpdate.MaxUnavailable.IntVal, nil
		}
		return deployment.Status.ReadyReplicas == deployment.Status.Replicas, nil
	})
}

func (s *StatefulSet) Wait(timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		set, err := s.Clientset().AppsV1().StatefulSets(s.Namespace).Get(s.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if set.Spec.UpdateStrategy.Type != appsv1.RollingUpdateStatefulSetStrategyType {
			return true, nil
		}

		var partition int
		var replicas = 1
		if set.Spec.UpdateStrategy.RollingUpdate != nil && set.Spec.UpdateStrategy.RollingUpdate.Partition != nil {
			partition = int(*set.Spec.UpdateStrategy.RollingUpdate.Partition)
		}
		if set.Spec.Replicas != nil {
			replicas = int(*set.Spec.Replicas)
		}

		expectedReplicas := replicas - partition
		if int(set.Status.UpdatedReplicas) != expectedReplicas {
			return false, nil
		}
		if int(set.Status.ReadyReplicas) != replicas {
			return false, nil
		}
		return true, nil
	})
}
