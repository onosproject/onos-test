## How To Use ONIT?

The primary interface for setting up test clusters and running tests is the `onit` command,
which provides a suite of commands for setting up and tearing down test clusters, adding
and removing [device simulators][simulators], adding and removing networks of [stratum] switches, adding and removing applications, running tests, and viewing test history. To see list of `onit commands` run `onit` from the shell as follows:
```bash
> onit 
Run onos integration tests on Kubernetes

Usage:
  onit [command]

Available Commands:
  add         Add resources to the cluster
  completion  Generated bash or zsh auto-completion script
  create      Create a test resource on Kubernetes
  debug       Open a debugger port to the given resource
  delete      Delete Kubernetes test resources
  fetch       Fetch resources from the cluster
  get         Get test configurations
  help        Help about any command
  onos-cli    Open onos-cli shell for executing commands
  remove      Remove resources from the cluster
  run         Run integration tests
  set         Set test configurations
  ssh         Open a ssh session to a node for executing remote commands

Flags:
  -h, --help   help for onit

Use "onit [command] --help" for more information about a command.
```

## Cluster Setup

The first step to running tests is to setup a test cluster with `onit create cluster`:

```bash
> onit create cluster
 ✓ Creating cluster namespace
 ✓ Setting up RBAC
 ✓ Setting up Atomix controller
 ✓ Starting Raft partitions
 ✓ Adding secrets
 ✓ Bootstrapping onos-topo cluster
 ✓ Bootstrapping onos-config cluster
 ✓ Setting up GUI
 ✓ Setting up CLI
 ✓ Creating ingress for services
cluster-8face0a8-bed6-11e9-a853-3c15c2cff232
```

You can also specify the number of nodes for each onos subsystem, for example, to create a cluster which runs 
two onos-config and two onos-topo pods, run the following command:
```bash
onit create cluster onit-1 --config-nodes 2 --topo-nodes 2
```

To setup the cluster, onit creates a unique namespace within which to create test resources,
deploys [Atomix][atomix] inside the test namespace, and configures and deploys onos-config nodes.
Once the cluster is setup, the command will output the name of the test namespace. The namespace
can be used to view test resources via `kubectl`:

```bash
> kubectl get pods -n cluster-8face0a8-bed6-11e9-a853-3c15c2cff232
NAME                                 READY   STATUS    RESTARTS   AGE
atomix-controller-7f95d69f47-sbsn4   1/1     Running   0          107s
onos-cli-6c6cf7cc89-fmvwq            1/1     Running   0          63s
onos-config-d68456bd7-xf9nv          1/1     Running   0          72s
onos-config-envoy-5c49b74dc4-vpmm4   1/1     Running   0          71s
onos-gui-6c78895d94-lj5vq            1/1     Running   0          65s
onos-topo-654999644-674z5            1/1     Running   0          87s
onos-topo-envoy-dcddf9dc6-2f5q9      1/1     Running   0          87s
raft-1-0                             1/1     Running   0          98s
```

The `create cluster` command supports additional flags for defining the cluster architecture:
```bash
  Flags:
    -c, --config string               test cluster configuration (default "default")
        --config-nodes int            the number of onos-config nodes to deploy (default 1)
        --docker-registry string      an optional host:port for a private Docker registry
    -h, --help                        help for cluster
        --image-pull-policy string    the Docker image pull policy (default "IfNotPresent")
        --image-tags stringToString   the image docker container tag for each node in the cluster (default [topo=debug,simulator=latest,stratum=latest,test=latest,atomix=latest,raft=latest,config=debug])
    -s, --partition-size int          the size of each Raft partition (default 1)
    -p, --partitions int              the number of Raft partitions to deploy (default 1)
        --topo-nodes int              the number of onos-topo nodes to deploy (default 1) 
```
Once the cluster is setup, the cluster configuration will be added to the `onit` configuration
and the deployed cluster will be set as the current cluster context:

```bash
> onit get cluster
cluster-b8c45834-a81c-11e9-82f4-3c15c2cff232
```

You can also create a cluster by passing a name to the `onit create cluster` command. To create a cluster with name `onit-2`, we run the following command:

```bash
> onit create cluster onit-1
onit create cluster onit-1
 ✓ Creating cluster namespace
 ✓ Setting up RBAC
 ✓ Setting up Atomix controller
 ✓ Starting Raft partitions
 ✓ Adding secrets
 ✓ Bootstrapping onos-topo cluster
 ✓ Bootstrapping onos-config cluster
 ✓ Setting up GUI
 ✓ Setting up CLI
 ✓ Creating ingress for services
onit-1
```

if we run `onit get clusters` command, we should be able to see the two clusters that we created:

```bash
> onit get clusters
ID                                             SIZE   PARTITIONS
cluster-b8c45834-a81c-11e9-82f4-3c15c2cff232   1      1
onit-1
```

When multiple clusters are deployed, you can switch between clusters by setting the current
cluster context:

```bash
> onit set cluster onit-1
onit-1
```

This will run all future cluster operations on the configured cluster. Alternatively, most commands support a flag to override the default cluster:

To delete a cluster, run `onit delete cluster`:
```bash
> onit delete cluster
✓ Deleting cluster namespace
```

## Adding Simulators

Most tests require devices to be added to the cluster. The `onit` command supports adding and
removing [device simulators][simulators] through the `add` and `remove` commands. To add a
simulator to the current cluster, simply run `onit add simulator`:

```bash
> onit add simulator 
✓ Setting up simulator
✓ Reconfiguring onos-config nodes
device-1186885096
```

When a simulator is added to the cluster, the cluster is reconfigured in two phases:
* Bootstrap a new [device simulator][simulators] with the provided configuration
* Reconfigure and redeploy the onos-config cluster with the new device in its stores

To give a name to a simulator, pass a name to `onit add simulator` command as follows
```bash
> onit add simulator sim-2
✓ Setting up simulator
✓ Reconfiguring onos-config nodes
sim-2
```

To get list of simulators, run `onit get simulators` as follows:

```bash
> onit get simulators 
device-1186885096
sim-2
```
Simulators can similarly be removed with the `remove simulator` command:

```bash
> onit remove simulator device-1186885096
 ✓ Tearing down simulator
 ✓ Reconfiguring onos-config nodes
```

As with the `add` command, removing a simulator requires that the onos-config cluster be reconfigured and redeployed.

## Adding Networks
To run some of the tests on stratum switches, we can create a network of stratum switches using Mininet. To create a network of stratum switches, we can use `onit add network [Name] [Mininet Options]` as follows: 

* To create a single node network, simply run `onit add network`. This command creates a single node network and assigns a name to it automatically. 
* To create a linear network topology with two switches and name it *stratum-linear*, simply run the following command:

```bash
> onit add network stratum-linear -- --topo linear,2
✓ Setting up network
✓ Reconfiguring onos-config nodes
stratum-linear
```

When a network is added to the cluster, the cluster is reconfigured in two phases:
* Bootstrap one or more than one new stratum switches with the provided configuration
* Reconfigure and redeploy the onos-config cluster with the new switches in its stores

To add a single node network topology, run the following command:
```bash
> onit add network
✓ Setting up network
✓ Reconfiguring onos-config nodes
network-2878434070
```

To get list of networks, run the following command:
```bash
> onit get networks
network-2878434070
stratum-linear
```
Networks can be removed using `onit remove network` command. For example, to remove the linear topolog that is created using the above command, you should run the following command:
```bash
> onit remove network stratum-linear
 ✓ Tearing down network
 ✓ Reconfiguring onos-config nodes
```

As with the `add` command, removing a network requires that the onos-config cluster be reconfigured and redeployed.

**Note**: In the current implementation, we support the following network topologies:

* A *Single* node network topology
* A *Linear* network topology

## Adding Applications

Applications from outside of `onit` can be added to an `onit` cluster using the `onit add app` command. This command takes as
the name of the app as an argument. It also has the `--image` flag that allows a user specify the image of 
the application that should be  deployed. The user also can specify the pull policy for the image using `--image-pull-policy` flag. For example, to deploy the latest version of `onos-ztp`  application:
```bash
> onit add app onos-ztp --image onosproject/onos-ztp:latest --image-pull-policy "Always" 
 ✓ Setting up app
 onos-ztp
```

To give a name to an app, pass a name to `onit add app` command as follows
```bash
> onit add app onosproject/onos-ztp:latest ztp
   ✓ Setting up app 
  ztp
```

To get list of apps, run `onit get apps` as follows:

```bash
> onit get apps
  app-128922186
```
Apps can be removed with the `remove app` command:

```bash
> onit remove app app-128922186
   ✓ Tearing down app 
```


## SSH Into A Cluster Node
onit allows you to ssh into a node using the following command:
```bash
onit ssh <name of a node>
```

onit also provides a command that you can run [onos-cli](onos-cli) commands via onit as follows:
```bash
onit onos-cli
~ $ onos
ONOS command line client

Usage:
  onos [command]

Available Commands:
  completion  Generated bash or zsh auto-completion script
  help        Help about any command
  topo
  ztp         ONOS zero-touch provisioning subsystem commands

Flags:
  -h, --help   help for onos

Use "onos [command] --help" for more information about a command.
```


[onos-cli]: http://github.com/onosproject/onos-cli
[simulators]: https://github.com/onosproject/simulators
[atomix]: https://github.com/atomix/atomix

