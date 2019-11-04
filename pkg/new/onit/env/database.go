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

package env

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// Database provides the database environment
type Database interface {
	// Partitions returns all database partitions
	Partitions(group string) []Partition

	// Partition returns the Partition for the given partition
	Partition(name string) Partition
}

var _ Database = &database{}

// database is an implementation of the Database interface
type database struct {
	*testEnv
}

func (e *database) Partitions(group string) []Partition {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"type": "database", "group": group}}
	list, err := e.atomixClient.K8sV1alpha1().Partitions(e.namespace).List(metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String()})
	if err != nil {
		panic(err)
	}

	partitions := make([]Partition, 0, len(list.Items))
	for _, partition := range list.Items {
		partitions = append(partitions, e.Partition(partition.Name))
	}
	return partitions
}

func (e database) Partition(name string) Partition {
	return &partition{
		service: newService(name, "database", e.testEnv),
	}
}
