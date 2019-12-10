# Getting started with ONIT

## Prerequisites

ONIT sets up test clusters and runs test jobs inside [Kubernetes]. For development and testing, we recommend using 
[KIND] or [MicroK8s], but ONIT can run on any Kubernetes cluster.

Additionally, [Golang] 1.12 or later is recommended for downloading/compiling the ONIT binary.

## Installation

To install ONIT, fetch and compile the ONIT binary using `go get`:

```bash
> GO111MODULE=on go get github.com/onosproject/onos-test/cmd/onit
```

The `onit` CLI supports auto-completion of commands for bash and zsh. To enable auto-completion:

* `bash` - Run `source $(onit completion bash)`
* `zsh` - Run `source <(onit completion zsh)`

You can optionally persist the output of the `onit completion` command to your shell profile.

## Usage

### Development Cluster

To setup a cluster for development, after having pushed all the new images to `kind` run:
```bash
onit create cluster --set onos-cli.enabled=true
```

### Integration Tests

To run a suite of tests, use the `onit run test` command, providing a test image to run:

```bash
onit test --image <test-image>
```

For example, to run onos-config suite tests using [kind] cluster:
```bash
git clone https://github.com/onosproject/onos-config.git
cd onos-config
make kind
onit test --image onosproject/onos-config-tests:latest --suite gnmi
```

Benchmarks can be run with the `onit run benchmark` command:

```bash
onit benchmark --image <test-image> 
```

To run a suite (e.g. a suite of onos-topo tests): 
```bash
onit test --image onosproject/onos-topo-tests:latest --suite topo
```

To run a specific test (e.g. an onos-config test):
```bash
onit test --image onosproject/onos-config-tests:latest --suite config --test TestTransaction
```

More usage examples are provided in [debugging](debugging.md) document.

[Kubernetes]: https://kubernetes.io/
[KIND]: https://github.com/kubernetes-sigs/kind
[MicroK8s]: https://microk8s.io/
[Golang]: https://golang.org/
