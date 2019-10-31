package setup2

// TestSetup is a test configuration
type TestSetup interface {
	// Cluster sets up the test cluster
	Cluster() ClusterSetup

	// Database sets up the test cluster database
	Database() DatabaseSetup

	// Topo sets up the topo cluster
	Topo() TopoSetup

	// Config sets up the config cluster
	Config() ConfigSetup

	// GUI sets up the test cluster GUI
	GUI() GUISetup

	// CLI sets up the test cluster CLI
	CLI() CLISetup
}

// TestRunner is a test setup runner
type TestRunner interface {
	TestSetup

	// SetUp sets up the test
	SetUp() error

	// Run runs the test
	Run() error

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

// DatabaseSetup configures the database
type DatabaseSetup interface {
	// SetImage sets the partition image
	SetImage(image string)

	// SetPartitions sets the number of database partitions
	SetPartitions(partitions int) DatabaseSetup

	// SetPartitionSize sets the number of nodes per partition
	SetPartitionSize(partitionSize int) DatabaseSetup
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
	// SetNodes sets the number of config nodes
	SetNodes(nodes int) ConfigSetup

	// SetImage sets the config service image
	SetImage(image string) ConfigSetup
}

// GUISetup configures the ONOS GUI
type GUISetup interface {
}

// CLISetup configures the ONOS CLI
type CLISetup interface {
	// SetImage sets the CLI image
	SetImage(image string) CLISetup
}
