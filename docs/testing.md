# Testing with ONIT

ONIT provides a rich API for setting up and operating on μONOS clusters in tests.

## Creating a test image

Tests are deployed as standalone [Docker] images. Each test image may contain one or more test or benchmark suites,
and test commands may filter test suites and tests via command line arguments.

To create a test image, create a main function for the tests and call the `onit.Main()` function:

```go
package main

import (
	"github.com/onosproject/onos-test/pkg/onit"
)

func main() {
	onit.Main()
}
```

The `onit.Main()` function will run tests based on the arguments provided via the `onit` CLI. If no test suites are 
registered (as in the example above), no tests will be run. To register test and benchmark suites, call
`onit.RegisterTests` or `onit.RegisterBenchmarks` respectively:

```go
package main

import (
	"github.com/onosproject/onos-test/pkg/onit"
	"github.com/onosproject/onos-test/test/atomix"
)

func main() {
	// Register Atomix tests
	onit.RegisterTests("atomix", &atomix.SmokeTestSuite{})
	onit.RegisterTests("atomix-ha", &atomix.HATestSuite{})
	
	// Register Atomix benchmarks
	onit.RegisterBenchmarks("atomix", &atomix.BenchmarkSuite{})

	onit.Main()
}
```

Each test suite should be assigned a unique name which can be used to filter suites in the same image when running
tests via the CLI:

```bash
> onit run test --image onosproject/onos-tests:latest --suite atomix-ha
```

## Writing a test suite

Tests are grouped into suites, and suites are defined by Golang structs. All tests within a suite are run sequentially 
within a shared environment, but multiple suites within a test image may be run in parallel, each in their own
Kubernetes namespace.

To create a test suite, extend `onit.TestSuite`:

```go
type MyTestSuite struct {
	onit.TestSuite
}
```

Tests in a suite must follow the `Test*` naming pattern and take a `*testing.T` as the sole argument:

```go
func (s *MyTestSuite) TestStuff(t *testing.T) {
	...
}
```

Test suites can also implement a variety of interfaces for setting up and tearing down the suite and individual tests:
* `SetupTestSuite` - Run once before all tests to set up the test suite
* `TearDownTestSuite` - Run once after all tests to tear down the test suite
* `SetupTest` - Run before each test method to set up the test
* `TearDownTest` - Run after each test method to tear down the test

## Setup API

When a test is deployed it is run in a namespace that is virtually empty. Test suites are responsible for setting up
the systems required by the tests. To assist test suites in setting up μONOS clusters, ONIT provides a _setup API_
that allows test suites to configure every component of the μONOS cluster.

The setup API is typically only used once in the setup method of a test suite:

```go
func (s *MyTestSuite) SetupTestSuite() {

}
```

### Setting up the cluster

To facilitate setting up test clusters, the setup API provides configuration interfaces for each μONOS subsystem:
* `Atomix` - Configure the Atomix controller
* `Database` - Configure database partitions used by μONOS services
* `Topo` - Configure the μONOS topology service
* `Config` - Configure the μONOS config service
* `CLI` - Configure the μONOS command line client

```go
func (s *MyTestSuite) SetupTestSuite() {
	setup.Database().
		SetPartitions(3).
		SetReplicasPerPartition(3)
	setup.Topo().
		SetReplicas(2)
	setup.Config().
		SetReplicas(1)
	...
}
```

Once the desired subsystems have been configured, set up the cluster by calling `SetupOrDie`:

```go
func (s *MyTestSuite) SetupTestSuite() {
	setup.Database().
		SetPartitions(3).
		SetReplicasPerPartition(3)
	setup.Topo().
		SetReplicas(2)
	setup.Config().
		SetReplicas(1)
	setup.SetupOrDie()
}
``` 

## Runtime API

The setup API only provides for the deployment of the basic services required by a μONOS cluster. Most tests also 
require network devices, applications, and other components. ONIT provides a runtime API called `env` for adding devices
to the cluster, deploying and redeploying applications, killing nodes, and much more.

To prevent tests from polluting each other, the `env` API should only be used within the test methods themselves. Each
test is responsible for adding and removing the resources it requires.

```go
func (s *MyTestSuite) TestDevices(t *testing.T) {
	device1 := env.NewSimulator().
		SetName("device-1").
		AddOrDie()
	device2 := env.NewSimulator().
		SetName("device-2").
		AddOrDie()
}
```

### Testing core services

The `env` API provides information about each of the services running within the μONOS test cluster. The APIs for each
service can be used to list nodes, kill nodes, execute commands on any node, and connect to northbound APIs for
relevant services.

To list the nodes in the topo service:

```go
nodes := env.Topo().Nodes()
```

To kill a node in the topo service:

```go
nodes[0].Kill()
```

To wait for the topo service to recover after a failure:

```go
env.Topo().AwaitReady()
```

To execute a command on the CLI:

```go
env.CLI().Execute("onos topo get devices")
```

To connect the the topo service's northbound API:

```go
conn, err := env.Topo().Connect()
```

### Managing applications

To add an application to the cluster, call `env.NewApp()`:

```go
func (s *MyTestSuite) TestZTP(t *testing.T) {
	ztp := env.NewApp().
		SetName("ztp").
		SetImage("onosproject/onos-ztp:latest").
		SetReplicas(2).
		AddOrDie()
}
```

Applications must be configured with an image and may specify the number of nodes to deploy. Additionally, application
configurations can be used to expose ports, add secrets, specify container arguments, and set other options:

* `SetReplicas` sets the number of replicas to deploy
* `SetImage` sets the image to deploy
* `SetPullPolicy` sets the image pull policy
* `AddPort` adds a named port to the application deployment. Note: the first port added to the application will be
used for health checking.
* `SetPorts` sets the named ports to expose for the application
* `SetDebug` sets whether to enable debug mode. When debug mode is enabled, containers will be deployed with the 
`SYS_PTRACE` ability enabled
* `AddSecret` adds a secret to the application deployment. The `path` is the path at which the secret will be mounted
inside the app pods. The `value` is the value of the secret to add.
* `SetSecrets` sets the secrets to attach to the application deployment
* `SetUser` overrides the user with which to run the application containers
* `SetPrivileged` sets whether to run the application in privileged mode
* `AddEnv` adds an environment variable to the application
* `SetEnv` sets the environment variables to pass to the containers
* `SetArgs` sets the arguments to pass to the containers

```go
func (s *MyTestSuite) TestZTP(t *testing.T) {
	ztp := env.NewApp().
		SetName("ztp").
		SetImage("onosproject/onos-ztp:latest").
		SetReplicas(2).
		AddPort("grpc", 5150).
		AddSecret("/certs/onf.cacrt", caCert).
		AddSecret("/certs/onos-ztp.cert", ztpCert).
		AddSecret("/certs/onos-ztp.key", ztpKey).
		SetArgs("-caPath=/certs/onf.cacrt", "-certPath=/certs/onos-ztp.cert", "-keyPath=/certs/onos-ztp.key").
		AddOrDie()
}
```

Once an application has been deployed, tests can query the application nodes, execute commands within the application's
pod, kill nodes, and more:

```go
func (s *MyTestSuite) TestZTP(t *testing.T) {
	ztp := env.NewApp().
		SetName("ztp").
		SetImage("onosproject/onos-ztp:latest").
		AddPort("grpc", 5150).
		SetUser(0).
		SetReplicas(2).
		AddOrDie()  
	
	// Execute a command on a ztp node to add a role
	ztp.Execute("onos ztp add role foo")
	
	// Kill a ztp node
	ztp.Nodes()[0].Kill()
}
```

### Managing device simulators

To add a device simulator to the cluster, use `env.NewSimulator()` and add the device with `AddOrDie()`:

```go
func (s *MyTestSuite) TestGNMI(t *testing.T) {
	device1 := env.NewSimulator().
		SetName("device-1").
		AddOrDie()
	device2 := env.NewSimulator().
		SetName("device-2").
		AddOrDie()
	
	doGnmiSet(device1)
	doGnmiSet(device2)
}
```

When a device simulator is added to the cluster, it will automatically be added to the μONOS topology.

To improve performance when setting up multiple device simulators, simulators can be deployed concurrently using 
`AddSimulators()`:

```go
func (s *MyTestSuite) TestGNMI(t *testing.T) {
	devices := env.AddSimulators().
		With(env.NewSimulator().SetName("device-1")).
		With(env.NewSimulator().SetName("device-2")).
		AddAllOrDie()
	
	doGnmiSet(device1)
	doGnmiSet(device2)
}
```

Simulators can be removed by simply calling `Remove()` or `RemoveOrDie()`:

```go
env.Simulator("device-1").RemoveOrDie()
```

### Managing Mininet networks

To add a Mininet network to the cluster, use `env.NewNetwork()`:

```go
func (s *MyTestSuite) TestNetwork(t *testing.T) {
	network := env.NewNetwork().
		SetTopo("linear,2", 2).
		AddOrDie()
}
```

When a network is added to the cluster, ONIT will create a service per device and add each device to the μONOS topology.
Devices can then be accessed via the network API:

```go
for _, device := range network.Devices() {
	doGnmiSet(device)
}
```

To remove a network and all its devices, simply call `Remove()` or `RemoveOrDie()`:

```go
network.RemoveOrDie()
```

[Docker]: https://www.docker.com/
[Atomix]: https://atomix.io
[stratum]: https://www.opennetworking.org/stratum/

