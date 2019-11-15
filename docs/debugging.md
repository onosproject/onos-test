# Debugging with ONIT

## Setting up a cluster

To set up a test cluster, use the `onit create cluster` command:

```bash
> onit create cluster
‣ 2019-11-15T10:45:21-08:00 onos Setup ONOS cluster
‣ 2019-11-15T10:45:21-08:00 onos Setup namespace
✓ 2019-11-15T10:45:21-08:00 onos Setup namespace
‣ 2019-11-15T10:45:21-08:00 onos Set up RBAC
✓ 2019-11-15T10:45:21-08:00 onos Set up RBAC
✓ 2019-11-15T10:45:21-08:00 onos Setup ONOS cluster
‣ 2019-11-15T10:45:21-08:00 onos Setup Atomix controller
✓ 2019-11-15T10:45:36-08:00 onos Setup Atomix controller
‣ 2019-11-15T10:45:36-08:00 onos Setup Raft partitions
‣ 2019-11-15T10:45:36-08:00 onos Setup onos-topo
‣ 2019-11-15T10:45:36-08:00 onos Setup onos-config
‣ 2019-11-15T10:45:36-08:00 onos Setup onos-cli service
✓ 2019-11-15T10:45:40-08:00 onos Setup onos-cli service
✓ 2019-11-15T10:45:45-08:00 onos Setup Raft partitions
✓ 2019-11-15T10:45:58-08:00 onos Setup onos-topo
✓ 2019-11-15T10:46:06-08:00 onos Setup onos-config
```

When a cluster is created, `onit` creates a new namespace in Kubernetes and deploys the default μONOS services in that
namespace. By default, ONIT will use the `onos` namespace, but the namespace can be overwritten by specifying a
cluster name:

```bash
> onit -c my-cluster create cluster
...
```

## Deploying a test network

The μONOS cluster is deployed with each of its subsystems but no devices or applications. The `onit add` and 
`onit remove` subcommands can be used to add and remove resources respectively.

To add a device simulator to the cluster, use the `onit add simulator` command:

```bash
> onit -c my-cluster add simulator -n device-1
‣ 2019-11-15T11:18:20-08:00 onos Add simulator driving-skink
✓ 2019-11-15T11:18:29-08:00 onos Add simulator driving-skink
```

Simulators can be removed with the `onit remove simulator` command:

```bash
> onit -c my-cluster remove simulator -n device-1
```

ONIT also supports Mininet networks. To add a Mininet network, use the `onit add network` command:

```bash
> onit add network -n network --topo linear,2 --devices 2
```

The `onit add network` command requires two flags:
* `--topo` defines the Mininet topology
* `--devices` informs onit of the number of devices created in the topology

By default, ONIT will deploy the `opennetworking/mn-stratum` image, but as with all subsystems the image can be
overridden by the `--image` flag.

## Managing applications

```bash
> onit add app -n ztp
```
