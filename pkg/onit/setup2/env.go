package setup2

// TestEnv provides the test environment to tests
type TestEnv interface {
	// GetPartitionNodes gets the nodes in a partition
	GetPartitionNodes(partition int) []string

	// GetTopoNodes gets a list of topo nodes
	GetTopoNodes() []string

	// GetConfigNodes gets a list of config nodes
	GetConfigNodes() []string

	// GetSimulators gets a list of simulators
	GetSimulators() []string

	// AddSimulator adds a simulator
	AddSimulator(opts ...SimulatorOption) error

	// RemoveSimulator removes a simulator
	RemoveSimulator(name string) error

	// GetApps gets a list of deployed apps
	GetApps() []string

	// GetAppNodes gets a list of app nodes
	GetAppNodes(name string) []string

	// AddApp adds an application
	AddApp(opts ...AppOption) error

	// RemoveApp removes an application
	RemoveApp(name string) error

	// KillNode kills a node
	KillNode(node string)
}

type appOptions struct {
	name  string
	nodes int
}

type AppOption interface {
	apply(opts *appOptions)
}

type appNameOption struct {
	name string
}

func (o appNameOption) apply(opts *appOptions) {
	opts.name = o.name
}

// WithAppName sets the app name
func WithAppName(name string) AppOption {
	return appNameOption{
		name: name,
	}
}

type appNodesOption struct {
	nodes int
}

func (o appNodesOption) apply(opts *appOptions) {
	opts.nodes = o.nodes
}

// WithAppNodes sets the app nodes
func WithAppNodes(nodes int) AppOption {
	return appNodesOption{
		nodes: nodes,
	}
}

type simulatorOptions struct {
	name string
}

type SimulatorOption interface {
	apply(opts *simulatorOptions)
}

type simulatorNameOption struct {
	name string
}

func (o simulatorNameOption) apply(opts *simulatorOptions) {
	opts.name = o.name
}

// WithSimulatorName sets the simulator name
func WithSimulatorName(name string) SimulatorOption {
	return simulatorNameOption{
		name: name,
	}
}
