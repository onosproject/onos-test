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

package setup

import "github.com/onosproject/onos-test/pkg/onit/cluster"

const kafkaRelease = "kafka"
const kafkaChart = "kafka"

// KafkaSetup is an interface for setting up Kafka
type KafkaSetup interface {
	// SetReplicas sets the number of replicas to deploy
	SetReplicas(replicas int) KafkaSetup

	// AddTopic adds a topic to the Kafka cluster
	AddTopic(name string, partitions, replicationFactor int) KafkaSetup
}

var _ KafkaSetup = &clusterKafkaSetup{}

func (s *clusterSetup) Kafka() KafkaSetup {
	chart := s.getChart(kafkaChart, func() *cluster.Chart {
		chart := s.cluster.Charts().New()
		chart.SetChart(kafkaChart)
		chart.SetRelease(kafkaRelease)
		chart.SetRepository("http://storage.googleapis.com/kubernetes-charts-incubator")
		chart.SetValue("replicas", 1)
		chart.SetValue("zookeeper.replicaCount", 1)
		chart.SetValue("configurationOverrides.\"log.message.timestamp.type\"", "LogAppendTime")
		return chart
	})
	return &clusterKafkaSetup{
		chart: chart,
	}
}

// clusterKafkaSetup is an implementation of the KafkaSetup interface
type clusterKafkaSetup struct {
	chart *cluster.Chart
}

func (s *clusterKafkaSetup) SetReplicas(replicas int) KafkaSetup {
	s.chart.SetValue("replicas", replicas)
	return s
}

func (s *clusterKafkaSetup) AddTopic(name string, partitions, replicationFactor int) KafkaSetup {
	s.chart.AddValue("topics", &KafkaTopic{
		Name:              name,
		Partitions:        partitions,
		ReplicationFactor: replicationFactor,
	})
	return s
}

func (s *clusterKafkaSetup) setup() error {
	return s.chart.Setup()
}

// KafkaTopic is a Kafka topic
type KafkaTopic struct {
	Name              string
	Partitions        int
	ReplicationFactor int
}
