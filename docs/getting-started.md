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

To run a suite of tests, use the `onit run test` command, providing a test image to run:

```bash
> onit test --image onosproject/onos-tests:latest
```

Benchmarks can be run with the `onit run benchmark` command:

```bash
> onit benchmark --image onosproject/onos-tests:latest
```

[Kubernetes]: https://kubernetes.io/
[KIND]: https://github.com/kubernetes-sigs/kind
[MicroK8s]: https://microk8s.io/
[Golang]: https://golang.org/
