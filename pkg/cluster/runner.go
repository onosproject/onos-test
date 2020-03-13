// Copyright 2019-present Open Networking Foundation.
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

package cluster

import (
	"bufio"
	"errors"
	"github.com/onosproject/onos-test/pkg/model"
	"os"
	"time"

	kube "github.com/onosproject/onos-test/pkg/kubernetes"
	"github.com/onosproject/onos-test/pkg/util/logging"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const namespace = "kube-test"

// Job is a job configuration
type Job struct {
	ID              string
	Image           string
	ImagePullPolicy corev1.PullPolicy
	ModelChecker    bool
	ModelData       map[string]string
	Args            []string
	Env             map[string]string
	Timeout         time.Duration
	Type            string
}

// NewRunner returns a new test job runner
func NewRunner() (*Runner, error) {
	return &Runner{
		client: kube.Namespace(namespace).Clientset(),
	}, nil
}

// Runner manages the test coordinator cluster
type Runner struct {
	client *kubernetes.Clientset
}

// Run runs the given job in the coordinator namespace
func (r *Runner) Run(job *Job) error {
	if err := r.ensureNamespace(); err != nil {
		return err
	}

	err := r.startJob(job)
	if err != nil {
		return err
	}

	step := logging.NewStep(job.ID, "Run job")
	step.Start()

	// Get the stream of logs for the pod
	pod, err := r.getPod(job, func(pod corev1.Pod) bool {
		return len(pod.Status.ContainerStatuses) > 0 &&
			pod.Status.ContainerStatuses[0].Ready
	})
	if err != nil {
		step.Fail(err)
		return err
	} else if pod == nil {
		return errors.New("cannot locate job pod")
	}

	req := r.client.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
		Container: "job",
		Follow:    true,
	})
	reader, err := req.Stream()
	if err != nil {
		step.Fail(err)
		return err
	}
	defer reader.Close()

	// Stream the logs to stdout
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logging.Print(scanner.Text())
	}

	// Get the exit message and code
	_, status, err := r.getStatus(job)
	if err != nil {
		step.Fail(err)
		return err
	}

	step.Complete()
	os.Exit(status)
	return nil
}

// ensureNamespace sets up the test namespace
func (r *Runner) ensureNamespace() error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := r.client.CoreV1().Namespaces().Create(ns)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return r.setupRBAC()
}

// setupRBAC sets up role based access controls for the cluster
func (r *Runner) setupRBAC() error {
	if err := r.createClusterRole(); err != nil {
		return err
	}
	if err := r.createClusterRoleBinding(); err != nil {
		return err
	}
	if err := r.createServiceAccount(); err != nil {
		return err
	}
	return nil
}

// createClusterRole creates the ClusterRole required by the Atomix controller and tests if not yet created
func (r *Runner) createClusterRole() error {
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespace,
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods",
					"pods/log",
					"pods/exec",
					"services",
					"endpoints",
					"persistentvolumeclaims",
					"events",
					"configmaps",
					"secrets",
					"serviceaccounts",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"namespaces",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"apps",
				},
				Resources: []string{
					"deployments",
					"daemonsets",
					"replicasets",
					"statefulsets",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"policy",
				},
				Resources: []string{
					"poddisruptionbudgets",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"batch",
				},
				Resources: []string{
					"jobs",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"rbac.authorization.k8s.io",
				},
				Resources: []string{
					"clusterroles",
					"clusterrolebindings",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"apiextensions.k8s.io",
				},
				Resources: []string{
					"customresourcedefinitions",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"cloud.atomix.io",
					"k8s.atomix.io",
				},
				Resources: []string{
					"*",
				},
				Verbs: []string{
					"*",
				},
			},
		},
	}
	_, err := r.client.RbacV1().ClusterRoles().Create(role)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createClusterRoleBinding creates the ClusterRoleBinding required by the test manager
func (r *Runner) createClusterRoleBinding() error {
	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespace,
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      namespace,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     namespace,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	_, err := r.client.RbacV1().ClusterRoleBindings().Create(roleBinding)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// createServiceAccount creates a ServiceAccount used by the test manager
func (r *Runner) createServiceAccount() error {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespace,
			Namespace: namespace,
		},
	}
	_, err := r.client.CoreV1().ServiceAccounts(namespace).Create(serviceAccount)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// startJob starts running a test job
func (r *Runner) startJob(job *Job) error {
	step := logging.NewStep(job.ID, "Starting job")
	step.Start()
	if err := r.createJob(job); err != nil {
		step.Fail(err)
		return err
	}
	if err := r.awaitJobRunning(job); err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// createJob creates the job to run tests
func (r *Runner) createJob(job *Job) error {
	step := logging.NewStep(job.ID, "Deploy job coordinator")
	step.Start()

	env := make([]corev1.EnvVar, 0, len(job.Env))
	for key, value := range job.Env {
		env = append(env, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}
	env = append(env, corev1.EnvVar{
		Name:  "SERVICE_NAMESPACE",
		Value: namespace,
	})
	env = append(env, corev1.EnvVar{
		Name:  "SERVICE_NAME",
		Value: job.ID,
	})
	env = append(env, corev1.EnvVar{
		Name: "POD_NAMESPACE",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	})
	env = append(env, corev1.EnvVar{
		Name: "POD_NAME",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	})

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: job.ID,
			Labels: map[string]string{
				"job":  job.ID,
				"type": job.Type,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"job": job.ID,
			},
			Ports: []corev1.ServicePort{
				{
					Name: "management",
					Port: 5000,
				},
			},
		},
	}
	if _, err := r.client.CoreV1().Services(namespace).Create(svc); err != nil {
		return err
	}

	var volumes []corev1.Volume
	var containers []corev1.Container
	if job.ModelChecker {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      job.ID,
				Namespace: namespace,
				Annotations: map[string]string{
					"job":  job.ID,
					"type": job.Type,
				},
			},
			Data: job.ModelData,
		}
		_, err := r.client.CoreV1().ConfigMaps(namespace).Create(cm)
		if err != nil {
			step.Fail(err)
			return err
		}

		volumes = []corev1.Volume{
			{
				Name: "models",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: job.ID,
						},
					},
				},
			},
			{
				Name: "data",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		}
		containers = []corev1.Container{
			{
				Name:            "job",
				Image:           job.Image,
				ImagePullPolicy: job.ImagePullPolicy,
				Args:            job.Args,
				Env:             env,
				Ports: []corev1.ContainerPort{
					{
						Name:          "management",
						ContainerPort: 5000,
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "data",
						MountPath: model.DataPath,
					},
				},
			},
			{
				Name:            "model-checker",
				Image:           "onosproject/model-checker:latest",
				ImagePullPolicy: job.ImagePullPolicy,
				Ports: []corev1.ContainerPort{
					{
						Name:          "model-checker",
						ContainerPort: model.CheckerPort,
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "data",
						MountPath: model.DataPath,
					},
					{
						Name:      "models",
						MountPath: model.ModelsPath,
						ReadOnly:  true,
					},
				},
			},
		}
	} else {
		containers = []corev1.Container{
			{
				Name:            "job",
				Image:           job.Image,
				ImagePullPolicy: job.ImagePullPolicy,
				Args:            job.Args,
				Env:             env,
				Ports: []corev1.ContainerPort{
					{
						Name:          "management",
						ContainerPort: 5000,
					},
				},
			},
		}
	}

	zero := int32(0)
	one := int32(1)
	batchJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      job.ID,
			Namespace: namespace,
			Annotations: map[string]string{
				"job":  job.ID,
				"type": job.Type,
			},
		},
		Spec: batchv1.JobSpec{
			Parallelism:  &one,
			Completions:  &one,
			BackoffLimit: &zero,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"job":  job.ID,
						"type": job.Type,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: namespace,
					RestartPolicy:      corev1.RestartPolicyNever,
					Containers:         containers,
					Volumes:            volumes,
				},
			},
		},
	}

	if job.Timeout > 0 {
		timeoutSeconds := int64(job.Timeout / time.Second)
		batchJob.Spec.ActiveDeadlineSeconds = &timeoutSeconds
	}

	_, err := r.client.BatchV1().Jobs(namespace).Create(batchJob)
	if err != nil {
		step.Fail(err)
		return err
	}
	step.Complete()
	return nil
}

// awaitJobRunning blocks until the test job creates a pod in the RUNNING state
func (r *Runner) awaitJobRunning(job *Job) error {
	for {
		pod, err := r.getPod(job, func(pod corev1.Pod) bool {
			return len(pod.Status.ContainerStatuses) > 0 &&
				pod.Status.ContainerStatuses[0].Ready
		})
		if err != nil {
			return err
		} else if pod != nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// getStatus gets the status message and exit code of the given pod
func (r *Runner) getStatus(job *Job) (string, int, error) {
	for {
		pod, err := r.getPod(job, func(pod corev1.Pod) bool {
			return len(pod.Status.ContainerStatuses) > 0 &&
				pod.Status.ContainerStatuses[0].State.Terminated != nil
		})
		if err != nil {
			return "", 0, err
		} else if pod != nil {
			state := pod.Status.ContainerStatuses[0].State
			if state.Terminated != nil {
				return state.Terminated.Message, int(state.Terminated.ExitCode), nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// getPod finds the Pod for the given test
func (r *Runner) getPod(job *Job, predicate func(pod corev1.Pod) bool) (*corev1.Pod, error) {
	pods, err := r.client.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: "job=" + job.ID,
	})
	if err != nil {
		return nil, err
	} else if len(pods.Items) > 0 {
		for _, pod := range pods.Items {
			if predicate(pod) {
				return &pod, nil
			}
		}
	}
	return nil, nil
}
