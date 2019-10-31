package k8s

import (
	"github.com/onosproject/onos-test/pkg/onit/setup2"
)

type testConfig struct {
	partitionImage string
	partitions     int
	partitionSize  int
	topoImage      string
	topoNodes      int
	configImage    string
	configNodes    int
}

// TestSetup is a k8s test setup
type TestSetup struct {
	cluster *ClusterSetup
	topo    *TopoSetup
	config  *ConfigSetup
}

func (s *TestSetup) Cluster() setup2.ClusterSetup {
	return s.cluster
}

func (s *TestSetup) Topo() setup2.TopoSetup {
	return s.topo
}

func (s *TestSetup) Config() setup2.ConfigSetup {
	return s.config
}

func (s *TestSetup) SetUp() error {
	if err := s.cluster.setup(); err != nil {
		return err
	}
	if err := s.cluster.waitForSetup(); err != nil {
		return err
	}
	if err := s.topo.setup(); err != nil {
		return err
	}
	if err := s.config.setup(); err != nil {
		return err
	}
	if err := s.topo.waitForSetup(); err != nil {
		return err
	}
	if err := s.config.waitForSetup(); err != nil {
		return err
	}
	return nil
}

func (s *TestSetup) TearDown() error {
	if err := s.cluster.teardown(); err != nil {
		return err
	}
	return nil
}

type ClusterSetup struct {
	config *testConfig
}

func (s *ClusterSetup) SetPartitions(partitions int) setup2.ClusterSetup {
	s.config.partitions = partitions
	return s
}

func (s *ClusterSetup) SetPartitionSize(partitionSize int) setup2.ClusterSetup {
	s.config.partitionSize = partitionSize
	return s
}

func (s *ClusterSetup) SetPartitionImage(image string) setup2.ClusterSetup {
	s.config.partitionImage = image
	return s
}

type TopoSetup struct {
	config *testConfig
}

func (s *TopoSetup) SetNodes(nodes int) setup2.TopoSetup {
	s.config.topoNodes = nodes
	return s
}

func (s *TopoSetup) SetImage(image string) setup2.TopoSetup {
	s.config.topoImage = image
	return s
}

type ConfigSetup struct {
	config *testConfig
}

func (s *ConfigSetup) SetNodes(nodes int) setup2.ConfigSetup {
	s.config.configNodes = nodes
	return s
}

func (s *ConfigSetup) SetImage(image string) setup2.ConfigSetup {
	s.config.configImage = image
	return s
}
