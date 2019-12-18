# Debugging with ONIT

## Setting up a cluster

To set up a test cluster, use the `onit create cluster` command:

```bash
> onit create cluster
‣ 2019-12-10T11:00:01-05:00 onit-better-dolphin Starting job
‣ 2019-12-10T11:00:01-05:00 onit-better-dolphin Deploy job coordinator
✓ 2019-12-10T11:00:01-05:00 onit-better-dolphin Deploy job coordinator
✓ 2019-12-10T11:00:03-05:00 onit-better-dolphin Starting job
‣ 2019-12-10T11:00:03-05:00 onit-better-dolphin Run job
‣ 2019-12-10T16:00:03Z onos Setup namespace
✓ 2019-12-10T16:00:03Z onos Setup namespace
‣ 2019-12-10T16:00:03Z onos Set up RBAC
✓ 2019-12-10T16:00:03Z onos Set up RBAC
‣ 2019-12-10T16:00:03Z onos Setup Atomix controller
✓ 2019-12-10T16:00:05Z onos Setup Atomix controller
‣ 2019-12-10T16:00:05Z onos Setup Raft partitions
‣ 2019-12-10T16:00:05Z onos Setup onos-topo
‣ 2019-12-10T16:00:05Z onos Setup onos-config
✓ 2019-12-10T16:00:13Z onos Setup Raft partitions
✓ 2019-12-10T16:00:23Z onos Setup onos-topo
✓ 2019-12-10T16:00:24Z onos Setup onos-config
✓ 2019-12-10T11:00:25-05:00 onit-better-dolphin Run job
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
‣ 2019-12-10T11:12:07-05:00 onit-optimal-snake Starting job
‣ 2019-12-10T11:12:07-05:00 onit-optimal-snake Deploy job coordinator
✓ 2019-12-10T11:12:07-05:00 onit-optimal-snake Deploy job coordinator
✓ 2019-12-10T11:12:08-05:00 onit-optimal-snake Starting job
‣ 2019-12-10T11:12:08-05:00 onit-optimal-snake Run job
‣ 2019-12-10T16:12:08Z onos Add simulator device-1
✓ 2019-12-10T16:12:11Z onos Add simulator device-1
✓ 2019-12-10T11:12:12-05:00 onit-optimal-snake Run job
```

Simulators can be removed with the `onit remove simulator` command:

```bash
> onit -c my-cluster remove simulator -n device-1
```

ONIT also supports Mininet networks. To add a Mininet network, use the `onit add network` command:

```bash
> onit add network -n my-network --topo linear,2 --devices 2
```

The `onit add network` command requires two flags:
* `--topo` defines the Mininet topology
* `--devices` informs onit of the number of devices created in the topology

By default, ONIT will deploy the `opennetworking/mn-stratum` image, but as with all subsystems the image can be
overridden by the `--image` flag.

## Managing applications

Applications are arbitrary deployments that can be added to the μONOS cluster via `onit`. To add an application,
use the `onit add app` command:

```bash
> onit -c my-cluster add app -n my-app --image onosproject/my-app:latest
```

When adding an application to the cluster you _must_ specify an image to deploy. The application may be assigned a
name, and if no name is assigned a random human-readable name will be generated. Additional flags are provided for
exposing ports, supplying and mounting secrets, configuring the environment, overriding container arguments, and more.

* `-r` `--replicas` sets the number of replicas to deploy
* `-i` `--image` sets the image to deploy
* `--image-pull-policy` sets the image pull policy
* `-p` `--port` is a mapping of named ports to expose in the application service. Example:
`onit add app ... -p grpc:5150 -p debug:40000`
* `-d` `--debug` enables debug mode for the application
* `-s` `--secret` is a mapping of secret paths and values to mount to the application pods. Keys indicate the absolute
path at which to mount each secret, and values may be either a local file to mount or a secret value. Example:
`onit add app ... --secret /credentials/password=rocks --secret /certs/tls.crt=./certs/my-app.crt --secret /certs/tls.key=./certs/my-app.key`
* `-u` `--user` overrides the user with which to run the application containers
* `--privileged` runs the application containers in privileged mode
* `-e` `--env` is a mapping of environment variables. Example:
`onit add app ... -e APP_HOST=0.0.0.0 -e APP_PORT=5150`
* Additional arguments are passed to the application containers. For example, to pass `-host=0.0.0.0 -port=5150` to
the application:
`onit add app -i myproject/my-app:latest -p 5150 -- -host=0.0.0.0 -port=5150`

These flags can be used to deploy secure applications with highly configurable environments. For example, to deploy
the [onos-ztp] application use the following command:

```bash
> onit add app -n onos-ztp -i onosproject/onos-ztp:latest -u 0 -p grpc=5150 -r 2 -s /certs/onf.cacrt=configs/certs/onf.cacrt -s /certs/onos-ztp.crt=configs/certs/service.crt -s /certs/onos-ztp.key=configs/certs/service.key -- -caPath=/certs/onf.cacrt -keyPath=/certs/onos-ztp.key -certPath=/certs/onos-ztp.crt
```

To remove an application, use the `onit remove app` command:

```bash
> onit -c my-cluster remove app -n onos-ztp
```

## Running the onos-gui
Running the µONOS GUI is achievable on top of an 'onit' created cluster, by loading
its Helm Chart in to that cluster. See [onos-gui] for more details.

[onos-ztp]: https://github.com/onosproject/onos-ztp
[onos-gui]: https://docs.onosproject.org/onos-gui/docs/deployment/
