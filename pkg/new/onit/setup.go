package onit

import (
	"github.com/atomix/atomix-k8s-controller/pkg/client/clientset/versioned"
	"github.com/onosproject/onos-test/pkg/new/kubetest"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// NewSetup returns a new onit Setup
func NewSetup(kube kubetest.KubeAPI) Setup {
	setup := &testSetup{
		namespace:    kube.Namespace(),
		kubeClient:   kubernetes.NewForConfigOrDie(kube.Config()),
		atomixClient: versioned.NewForConfigOrDie(kube.Config()),
	}
	setup.atomix = &atomixSetup{
		testSetup: setup,
	}
	setup.database = &databaseSetup{
		testSetup: setup,
	}
	setup.topo = &topoSetup{
		testSetup: setup,
	}
	setup.config = &configSetup{
		testSetup: setup,
	}
	return setup
}

// Setup is an interface for setting up ONOS clusters
type Setup interface {
	Atomix() AtomixSetup
	Database() DatabaseSetup
	Topo() TopoSetup
	Config() ConfigSetup
	Setup()
}

// AtomixSetup is an interface for setting up the Atomix controller
type AtomixSetup interface {
	Setup
	Image(image string) AtomixSetup
}

// DatabaseSetup is an interface for setting up Raft partitions
type DatabaseSetup interface {
	Setup
	Image(image string) DatabaseSetup
	PullPolicy(pullPolicy corev1.PullPolicy) DatabaseSetup
	Partitions(partitions int) DatabaseSetup
	Nodes(nodes int) DatabaseSetup
}

// TopoSetup is an interface for setting up topo nodes
type TopoSetup interface {
	Setup
	Image(image string) TopoSetup
	PullPolicy(pullPolicy corev1.PullPolicy) TopoSetup
	Nodes(nodes int) TopoSetup
}

// ConfigSetup is an interface for setting up config nodes
type ConfigSetup interface {
	Setup
	Image(image string) ConfigSetup
	PullPolicy(pullPolicy corev1.PullPolicy) ConfigSetup
	Nodes(nodes int) ConfigSetup
}

// SimulatorSetup is an interface for setting up a simulator
type SimulatorSetup interface {
	Image(image string) SimulatorSetup
	PullPolicy(pullPolicy corev1.PullPolicy) SimulatorSetup
}

// NetworkSetup is an interface for setting up a network
type NetworkSetup interface {
	Image(image string) NetworkSetup
	PullPolicy(pullPolicy corev1.PullPolicy) NetworkSetup
}

// testSetup is an implementation of the Setup interface
type testSetup struct {
	namespace    string
	kubeClient   *kubernetes.Clientset
	atomixClient *versioned.Clientset
	atomix       *atomixSetup
	database     *databaseSetup
	topo         *topoSetup
	config       *configSetup
}

func (s *testSetup) Atomix() AtomixSetup {
	return s.atomix
}

func (s *testSetup) Database() DatabaseSetup {
	return s.database
}

func (s *testSetup) Topo() TopoSetup {
	return s.topo
}

func (s *testSetup) Config() ConfigSetup {
	return s.config
}

func (s *testSetup) Setup() {

}

var _ Setup = &testSetup{}

// atomixSetup is an implementation of the AtomixSetup interface
type atomixSetup struct {
	*testSetup
	image      string
	pullPolicy corev1.PullPolicy
}

func (s *atomixSetup) Image(image string) AtomixSetup {
	s.image = image
	return s
}

func (s *atomixSetup) PullPolicy(pullPolicy corev1.PullPolicy) AtomixSetup {
	s.pullPolicy = pullPolicy
	return s
}

var _ AtomixSetup = &atomixSetup{}

// databaseSetup is an implementation of the DatabaseSetup interface
type databaseSetup struct {
	*testSetup
	image      string
	pullPolicy corev1.PullPolicy
	partitions int
	nodes      int
}

func (s *databaseSetup) Image(image string) DatabaseSetup {
	s.image = image
	return s
}

func (s *databaseSetup) PullPolicy(pullPolicy corev1.PullPolicy) DatabaseSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *databaseSetup) Partitions(partitions int) DatabaseSetup {
	s.partitions = partitions
	return s
}

func (s *databaseSetup) Nodes(nodes int) DatabaseSetup {
	s.nodes = nodes
	return s
}

var _ DatabaseSetup = &databaseSetup{}

// topoSetup is an implementation of the TopoSetup interface
type topoSetup struct {
	*testSetup
	image      string
	pullPolicy corev1.PullPolicy
	nodes      int
}

func (s *topoSetup) Image(image string) TopoSetup {
	s.image = image
	return s
}

func (s *topoSetup) PullPolicy(pullPolicy corev1.PullPolicy) TopoSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *topoSetup) Nodes(nodes int) TopoSetup {
	s.nodes = nodes
	return s
}

var _ TopoSetup = &topoSetup{}

// configSetup is an implementation of the ConfigSetup interface
type configSetup struct {
	*testSetup
	image      string
	pullPolicy corev1.PullPolicy
	nodes      int
}

func (s *configSetup) Image(image string) ConfigSetup {
	s.image = image
	return s
}

func (s *configSetup) PullPolicy(pullPolicy corev1.PullPolicy) ConfigSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *configSetup) Nodes(nodes int) ConfigSetup {
	s.nodes = nodes
	return s
}

var _ ConfigSetup = &configSetup{}

// simulatorSetup is an implementation of the SimulatorSetup interface
type simulatorSetup struct {
	*testEnv
	image      string
	pullPolicy corev1.PullPolicy
}

func (s *simulatorSetup) Image(image string) SimulatorSetup {
	s.image = image
	return s
}

func (s *simulatorSetup) PullPolicy(pullPolicy corev1.PullPolicy) SimulatorSetup {
	s.pullPolicy = pullPolicy
	return s
}

var _ SimulatorSetup = &simulatorSetup{}

// networkSetup is an implementation of the NetworkSetup interface
type networkSetup struct {
	*testEnv
	image      string
	pullPolicy corev1.PullPolicy
}

func (s *networkSetup) Image(image string) NetworkSetup {
	s.image = image
	return s
}

func (s *networkSetup) PullPolicy(pullPolicy corev1.PullPolicy) NetworkSetup {
	s.pullPolicy = pullPolicy
	return s
}

var _ NetworkSetup = &networkSetup{}
