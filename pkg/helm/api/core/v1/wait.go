package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

func (p *Pod) Wait(timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		pod, err := p.Clientset().CoreV1().Pods(p.Namespace).Get(p.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, c := range pod.Status.Conditions {
			if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})
}

func (s *Service) Wait(timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		service, err := s.Clientset().CoreV1().Services(s.Namespace).Get(s.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if service.Spec.Type == corev1.ServiceTypeExternalName {
			return true, nil
		}
		if service.Spec.ClusterIP == "" {
			return false, nil
		}
		if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
			if len(service.Spec.ExternalIPs) > 0 {
				return true, nil
			}
			if service.Status.LoadBalancer.Ingress == nil {
				return false, nil
			}
		}
		return true, nil
	})
}
