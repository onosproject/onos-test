# How To Install and Run ONIT?

## Setup
The integration test framework is designed to operate on a Kubernetes cluster. It's recommended
that users use a local Kubernetes cluster suitable for development, e.g. [Minikube], [kind],
or [MicroK8s].

### Configuration

The test framework is controlled through the `onit` command. To install the `onit` command,
use `go get`:

```bash
> go get github.com/onosproject/onos-test/cmd/onit
```

To interact with a Kubernetes cluster, the `onit` command must have access to a local
Kubernetes configuration. Onit expects the same configuration as `kubectl` and will connect
to the same Kubernetes cluster as `kubectl` will connect to, so to determine which Kubernetes 
cluster onit will use, simply run `kubectl cluster-info`:

```bash
> kubectl cluster-info
Kubernetes master is running at https://127.0.0.1:49760
KubeDNS is running at https://127.0.0.1:49760/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
```

See the [Kubernetes documentation](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/)
for details on configuring both `kubectl` and `onit` to connect to a new cluster or multiple
clusters.

The `onit` command also maintains some cluster metadata in a local configuration file. The search
path for the `onit.yaml` configuration file is:
* `~/.onos`
* `/etc/onos`
* `.`

Users do not typically need to modify the `onit.yaml` configuration file directly. The onit 
configuration is primarily managed through various onit commands like `onit set`, `onit create`,
`onit add`, etc. It's recommended that users avoid modifying the onit configuration
file, but it's nevertheless important to note that the application must have write access to one
of the above paths.


### Onit Auto-Completion
*Onit* supports shell auto-completion for its various commands, sub-commands and flags.
You can enable this feature for *bash* or *zsh* as follows:
#### Bash Auto-Completion
To enable this for **bash**, run the following from the shell:

```bash
> eval "$(onit completion bash)"
```
#### Zsh Auto-Completion 

To enable this for **zsh**, run the following from the shell:
```bash
> source <(onit completion zsh)
```

**Note**: We also recommend to add the output of the above commands to *.bashrc* or *.zshrc*.

### Docker

The `onit` command manages clusters and runs tests by deploying locally built [Docker] containers
on [Kubernetes]. Docker image builds are an essential component of the `onit` workflow. **Each time a
change is made to either the core or integration tests, Docker images must be rebuilt** and made
available within the Kubernetes cluster in which tests are being run. The precise process for building
Docker images and adding them to a local Kubernetes cluster is different for each setup.

#### Building for Minikube

[Minikube] runs a VM with its own Docker daemon running inside it. To build the Docker images
for Minikube, ensure you use configure your shell to use Minikube Docker context before building:

```bash
> eval $(minikube docker-env)
```

Once the shell has been configured, use `make images` to build the Docker images:

```bash
> make images
```

Note that `make images` _must be run every time a change is made_ to either the core code
or integration tests.

#### Building for Kind

[Kind][kind] provides an alternative to [Minikube] which runs Kubernetes in a Docker container.

Assuming you have dowloaded kind as per [instructions][kind-install], the first time you boot the kind cluster 
or if you have rebooted your docker deamon you need to issue:

```bash
> kind create cluster
```

and for each window you intend to use `onit` commands in you will need to export the `KUBECONFIG` 
variable:

```bash
> export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
```

As with Minikube, kind requires specific setup to ensure Docker images modified and built
locally can be run within the kind cluster. Rather than switching your Docker environment to
a remote Docker server, kind requires that images be explicitly loaded into the cluster each
time they're built. For this reason, we provide a convenience make target: `kind`:

```bash
> make kind
```

When the `kind` target is run, the `onos-config` and `onos-config-integration-tests` images will
be built and loaded into the kind cluster, so no additional step is necessary.

