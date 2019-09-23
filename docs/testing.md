# Running Tests

## Running single Tests

Once the cluster has been setup for the test, to run a test simply use `onit run`:

```bash
> onit run test single-path
 ✓ Starting test job: test-25324770
=== RUN   single-path
--- PASS: single-path (0.46s)
PASS
PASS
```

You can specify as many tests as desired:

```bash
> onit run test single-path transaction subscribe
...
```

## Running a suite of Tests

`onit` can also run a suite of tests e.g. `integration-tests` which encompasses all the active integration tests.
```bash
> onit run suite integration-tests
 ✓ Starting test job: test-3109317976
=== RUN   single-path
--- PASS: single-path (0.20s)
=== RUN   subscribe
--- PASS: subscribe (0.09s)
PASS
```

## Test Run logs

Each test run is recorded as a job in the Kubernetes cluster. This ensures that logs, statuses,
and exit codes are retained for the lifetime of the cluster. Onit supports viewing past test
runs and logs via the `get` command:

```bash
> onit get history
ID                TESTS                     STATUS   EXIT CODE   MESSAGE
test-25324770     test,single-path          PASSED   0
test-2886892866   test,subscribe            PASSED   0
test-3109317976   suite,integration-tests   PASSED   0
```

To get the logs from a specific test, use `onit get logs` with the test ID:

```bash
> onit get logs test-2886892866
=== RUN   test-single-path-test
--- PASS: test-single-path-test (0.04s)
PASS
```

## Debugging

The `onit` command provides a set of commands for debugging test clusters. The `onit` command
can be used to `get logs` for every resource deployed in the test cluster. Simply pass the
resource ID (e.g. test `ID`, node `ID`, partition `ID`, etc) to the `onit get logs` command
to get the logs for a resource.

To list all types of nodes (e.g. onos-topo, onos-config, etc) running in the cluster, use `onit get nodes`, the output will be like the following:

```bash
> onit get nodes
onit get nodes
ID                             TYPE     STATUS
onos-topo-7cd788fb7f-2zvsp     topo     RUNNING
onos-topo-7cd788fb7f-rc6m5     topo     RUNNING
onos-config-6f8fcf5954-55zn2   config   RUNNING
onos-config-6f8fcf5954-pglkz   config   RUNNING
```


To get logs for the above node, run the following command:
```bash
> onit get logs onos-config-569c7d8546-jscg8
I0625 21:55:32.027255       1 onos-config.go:114] Starting onos-config
I0625 21:55:32.030184       1 manager.go:98] Configuration store loaded from /etc/onos-config/configs/configStore.json
I0625 21:55:32.030358       1 manager.go:105] Change store loaded from /etc/onos-config/configs/changeStore.json
I0625 21:55:32.031087       1 manager.go:112] Device store loaded from /etc/onos-config/configs/deviceStore.json
I0625 21:55:32.031222       1 manager.go:119] Network store loaded from /etc/onos-config/configs/networkStore.json
I0625 21:55:32.031301       1 manager.go:47] Creating Manager
...
```

To list the Raft partitions running in the cluster, use `onit get partitions`:

```bash
> onit get partitions
ID   GROUP   NODES
1    raft    raft-1-0
```
To get logs for the above partions, run the following command:
```bash
> onit get logs raft-1-0
21:10:24.466 [main] INFO  io.atomix.server.AtomixServerRunner - Node ID: raft-1-0
21:10:24.472 [main] INFO  io.atomix.server.AtomixServerRunner - Partition Config: /etc/atomix/partition.json
21:10:24.472 [main] INFO  io.atomix.server.AtomixServerRunner - Protocol Config: /etc/atomix/protocol.json
21:10:24.473 [main] INFO  io.atomix.server.AtomixServerRunner - Starting server
...
```

To list the tests that have been run, use `onit get history`:

```bash
> onit get history
ID                                     TESTS                   STATUS   EXIT CODE   MESSAGE
3cf7311a-9776-11e9-bfc3-acde48001122   test-integration-test   PASSED   0
68ad9154-977c-11e9-bcf2-acde48001122   test-integration-test   FAILED   1
71a0623c-977c-11e9-8478-acde48001122   test-single-path-test   PASSED   0
9e512cdc-9720-11e9-ba6e-acde48001122   *                       PASSED   0
da629d06-9774-11e9-bb50-acde48001122   *                       PASSED   0
```
To get logs for one of the above histories, run the following command:
```bash
> onit get logs 71a0623c-977c-11e9-8478-acde48001122
=== RUN   test-single-path-test
--- PASS: test-single-path-test (0.04s)
PASS
```

To download logs from a node, you can run `onit fetch logs` command. For example, to download logs from *onos-config-66d54956f5-xwpsh* node, run the following command:
```bash
onit fetch logs onos-config-66d54956f5-xwpsh
```

You can refer to [Debug onos-config in Onit Using Delve](../../onos-config/docs/debugging.md) to learn more about debugging of onos-config pod using [*Delve*](https://github.com/go-delve/delve) debugger.


## API

Tests are implemented using Go's `testing` package;

```go
func MyTest(t *testing.T) {
	t.Fail("you messed up!")
}
```

However, rather than running tests using `go test`, we provide a custom registry of tests to
allow human-readable names to be assigned to tests for ease of use. Once you've written a test,
register the test in an `init` function:

```go
func init() {
	Registry.Register("my-test", MyTest)
}
```

Once a test has been registered, you should be able to see the test via the `onit` command:

```bash
> onit get tests
my-test
...
```

The test framework also provides the capability of adding your test to a suite defined in `suites.go`.
To see the suites you can execute:
```bash
> onit get suites
SUITE               TESTS
alltests            single-path, subscribe, transaction
sometests           subscribe, transaction
integration-tests   single-path
```

To add your test to a suite in the init function the register method must be called with the suites parameter:
```go
func init() {
    Registry.RegisterTest("my-test", MyTest, []*runner.TestSuite{AllTests})
}
```

The test framework provides utility functions for creating clients and other resources within
the test environment. The test environment is provided by the `env` package:

```go
client, err := env.NewGnmiClient(context.Background(), "")
...
```

When devices are deployed in the test configuration, a list of device IDs can be retrieved from
the environment:

```go
devices := env.GetDevices()
```

[Kubernetes]: https://kubernetes.io
[Minikube]: https://kubernetes.io/docs/setup/learning-environment/minikube/
[kind]: https://github.com/kubernetes-sigs/kind
[kind-install]: https://github.com/kubernetes-sigs/kind#installation-and-usage
[MicroK8s]: https://microk8s.io/
[Docker]: https://www.docker.com/
[Atomix]: https://atomix.io
[simulators]: https://github.com/onosproject/simulators
[stratum]: https://www.opennetworking.org/stratum/
[onos](https://github.com/onosproject)
[onos-test](https://github.com/onosproject/onos-test)
