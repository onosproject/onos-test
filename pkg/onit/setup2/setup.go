package setup2

// TestSetup is a test configuration
type TestSetup interface {
	// Cluster sets up the test cluster
	Cluster() ClusterSetup

	// Topo sets up the topo cluster
	Topo() TopoSetup

	// Config sets up the config cluster
	Config() ConfigSetup

	// SetUp sets up the test
	SetUp() error

	// TearDown tears down the test
	TearDown() error
}

// ClusterSetup configures the cluster
type ClusterSetup interface {
	// SetPartitions sets the partitions
	SetPartitions(partitions int) ClusterSetup

	// SetPartitionSize sets the partition size
	SetPartitionSize(partitionSize int) ClusterSetup

	// SetPartitionImage sets the partition image name
	SetPartitionImage(image string) ClusterSetup
}

// TopoSetup configures the topology service
type TopoSetup interface {
	// SetNodes sets the number of topo nodes
	SetNodes(nodes int) TopoSetup

	// SetImage sets the topo service image
	SetImage(image string) TopoSetup
}

// ConfigSetup configures the config service
type ConfigSetup interface {
	TestSetup

	// SetNodes sets the number of config nodes
	SetNodes(nodes int) ConfigSetup

	// SetImage sets the config service image
	SetImage(image string) ConfigSetup
}
