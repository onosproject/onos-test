# How To Install and Run ONIT?

## Setup
The integration test framework is designed to operate on a Kubernetes cluster. It's recommended
that users use a local Kubernetes cluster suitable for development, e.g. [k3d], [Minikube], [kind],
or [MicroK8s]. To run `onit`, you need to install `go` tools on your machine as explained [here](https://golang.org/doc/install)

### Configuration

The test framework is controlled through the `onit` command. To install the `onit` command,
use `go get`:

```bash
export GO111MODULE=on
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



### Build Using K8S Local Cluster Tools


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

[Kind] provides an alternative to [Minikube] which runs Kubernetes in a Docker container.

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

if you run the above command from the root of onos-test,  
the `onos-test-runner` image will
be built and loaded into the kind cluster. 

**Note:** The same make target (i.e. make kind) is provided
in other onos subsystems Makefiles such as [onos-config], [onos-topo], etc that allows you to
build and load other onos subsystem docker images into the kind cluster.

#### Building for MicroK8s
[microk8s](https://microk8s.io/) is a Kubernetes cluster solution that runs on Ubuntu
and other platforms. On Ubuntu is installed through the **snap** system on Ubuntu
16 and above.

After installing with:
```bash
snap install microk8s --classic
```
install **kubectl**. This can also be installed with snap:
```bash
sudo snap install kubectl --classic
```

Add your user name to the microk8s group in Linux:
```bash
usermod -a -G microk8s $USER
```
It will be necessary to log out and back in again to enable this.

For **onit** to work the **dns** service at least has to be installed. It is also
convenient to also install the dashboard.
```bash
microk8s.enable dns,dashboard
```

Also when **onit** runs it needs inter pod communication. Depending on the
installation of on Ubuntu the firewall may need to be disabled. On Ubuntu, the
system must be rebooted first and the **cbr0** interface should be visible.

Then run:
```bash
sudo ufw allow in on cbr0 && sudo ufw allow out on cbr0
sudo ufw default allow routed
```

Running
```bash
microk8s.inspect
```
should show no warnings about firewall.

After this **onit** can start the cluster. On Ubuntu running the debug versions
of **onos-topo** and **onos-config** has a problem in starting because of the
need for root permissions. For this reason **onit** __must__ be started as the 
**latest** version with:
```bash
onit create cluster --image-tags="topo=latest,config=latest"
```

By default MicroK8s will pull docker images from docker hub, and not the local
machine. To load a local image in to Microk8s:

```bash
docker save mynginx > myimage.tar
microk8s.ctr -n k8s.io image import myimage.tar
```

(see <https://microk8s.io/docs/working> for more details).

[Kubernetes]: https://kubernetes.io
[Minikube]: https://kubernetes.io/docs/setup/learning-environment/minikube/
[kind]: https://github.com/kubernetes-sigs/kind
[kind-install]: https://github.com/kubernetes-sigs/kind#installation-and-usage
[MicroK8s]: https://microk8s.io/
[k3d]: https://github.com/rancher/k3d
[kind-install]: https://github.com/kubernetes-sigs/kind#installation-and-usage
[onos-test]: https://github.com/onosproject/onos-test
[onos-config]: https://github.com/onosproject/onos-config
[onos-topo]: https://github.com/onosproject/onos-topo