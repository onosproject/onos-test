package onit

import (
	"github.com/atomix/atomix-k8s-controller/pkg/client/clientset/versioned"
	"github.com/onosproject/onos-test/pkg/new/kubetest"
	"k8s.io/client-go/kubernetes"
)

// NewEnv returns a new onit Setup
func NewEnv(kube kubetest.KubeAPI) Env {
	env := &testEnv{
		namespace:    kube.Namespace(),
		kubeClient:   kubernetes.NewForConfigOrDie(kube.Config()),
		atomixClient: versioned.NewForConfigOrDie(kube.Config()),
	}
	env.atomix = &atomixEnv{
		testEnv: env,
	}
	env.database = &databaseEnv{
		testEnv: env,
	}
	env.topo = &topoEnv{
		serviceEnv: &serviceEnv{
			testEnv: env,
		},
	}
	env.config = &configEnv{
		serviceEnv: &serviceEnv{
			testEnv: env,
		},
	}
	env.simulators = &simulatorsEnv{
		serviceEnv: &serviceEnv{
			testEnv: env,
		},
	}
	env.networks = &networksEnv{
		serviceEnv: &serviceEnv{
			testEnv: env,
		},
	}
	return env
}

// Env is an interface for tests to operate on the ONOS environment
type Env interface {
	// Atomix returns the Atomix environment
	Atomix() AtomixEnv

	// Database returns the database environment
	Database() DatabaseEnv

	// Topo returns the topo environment
	Topo() TopoEnv

	// Config returns the config environment
	Config() ConfigEnv

	// Simulators returns the simulators environment
	Simulators() SimulatorsEnv

	// Networks returns the networks environment
	Networks() NetworksEnv
}

// AtomixEnv provides the Atomix environment
type AtomixEnv interface {
	// Nodes returns the Atomix controller nodes
	Nodes() []string
}

// DatabaseEnv provides the database environment
type DatabaseEnv interface {
	// Partitions returns the number of database partitions
	Partitions() int

	// Nodes returns the nodes in the given partition
	Nodes(partition int) []string
}

// ServiceEnv is a base interface for service environments
type ServiceEnv interface {
	// Nodes returns the service nodes
	Nodes() []string

	// Kill kills the given node
	Kill(node string)
}

// TopoEnv provides the topo environment
type TopoEnv interface {
	ServiceEnv
}

// ConfigEnv provides the config environment
type ConfigEnv interface {
	ServiceEnv
}

// SimulatorsEnv provides the simulators environment
type SimulatorsEnv interface {
	ServiceEnv
}

// NetworksEnv provides the networks environment
type NetworksEnv interface {
	ServiceEnv
}

// testEnv is an implementation of the Env interface
type testEnv struct {
	namespace    string
	kubeClient   *kubernetes.Clientset
	atomixClient *versioned.Clientset
	atomix       *atomixEnv
	database     *databaseEnv
	topo         *topoEnv
	config       *configEnv
	simulators   *simulatorsEnv
	networks     *networksEnv
}

func (e *testEnv) Atomix() AtomixEnv {
	return e.atomix
}

func (e *testEnv) Database() DatabaseEnv {
	return e.database
}

func (e *testEnv) Topo() TopoEnv {
	return e.topo
}

func (e *testEnv) Config() ConfigEnv {
	return e.config
}

func (e *testEnv) Simulators() SimulatorsEnv {
	return e.simulators
}

func (e *testEnv) Networks() NetworksEnv {
	return e.networks
}

// atomixEnv is an implementation of the AtomixEnv interface
type atomixEnv struct {
	*testEnv
}

func (e *atomixEnv) Nodes() []string {
	panic("implement me")
}

var _ AtomixEnv = &atomixEnv{}

// databaseEnv is an implementation of the DatabaseEnv interface
type databaseEnv struct {
	*testEnv
}

func (e *databaseEnv) Partitions() int {
	panic("implement me")
}

func (e *databaseEnv) Nodes(partition int) []string {
	panic("implement me")
}

var _ DatabaseEnv = &databaseEnv{}

// serviceEnv is an implementation of the ServiceEnv interface
type serviceEnv struct {
	*testEnv
}

func (e *serviceEnv) Nodes() []string {
	panic("implement me")
}

func (e *serviceEnv) Kill(node string) {
	panic("implement me")
}

var _ ServiceEnv = &serviceEnv{}

// topoEnv is an implementation of the TopoEnv interface
type topoEnv struct {
	*serviceEnv
}

var _ TopoEnv = &topoEnv{}

// configEnv is an implementation of the ConfigEnv interface
type configEnv struct {
	*serviceEnv
}

var _ ConfigEnv = &configEnv{}

// simulatorsEnv is an implementation of the SimulatorsEnv interface
type simulatorsEnv struct {
	*serviceEnv
}

var _ SimulatorsEnv = &simulatorsEnv{}

// networksEnv is an implementation of the NetworksEnv interface
type networksEnv struct {
	*serviceEnv
}

var _ NetworksEnv = &networksEnv{}
