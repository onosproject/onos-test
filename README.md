<!--
SPDX-FileCopyrightText: 2022-present Intel Corporation
SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
SPDX-License-Identifier: Apache-2.0
-->

## onos-test : µONOS Architecture Integration and Smoke Tests
[![Go Report Card](https://goreportcard.com/badge/github.com/onosproject/onos-test)](https://goreportcard.com/report/github.com/onosproject/onos-test)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/onosproject/onos-test?status.svg)](https://godoc.org/github.com/onosproject/onos-test)

This repository contains testing infrastructure used for the [onos project] as well as testing utilities. The tests 
are run as part of Continuous Integration (CI) testing on Jenkins and by developers testing their code.

# Integration Tests
Each component of the µONOS project implements its own integration tests. These tests are written in Go and use
[helmit] for running end-to-end integration tests on [Kubernetes]. You can see some examples of helmit based integration
tests in [onos config], [onos e2t], and [onos MLB].

Integration tests are run in one of 3 ways (onos-config is used here for examples):
* From the onos-test repository, this is mostly used by Jenkins CI. This invocation will spin up a [kind] cluster to run the tests:

  `make onos-config-integration-tests`
* From the component repository using make:

  `cd onosproject/onos-config; make integration-tests`

* Manually on the command line:

  `cd onosproject/onos-config; kubectl create ns test; helmit test -n test ./cmd/onos-config-tests --suite config`
  
# Smoke tests
Smoke tests are scripts written in bash and python that exercise multiple components working together in a cluster.
These tests load the ONOS charts into a [kind] instance and use the [onos CLI] to interrogate the components to be
sure they are operating properly. Smoke tests are intended to run either as part of Jenkins CI or by developers to test changes.

The behavior of smoke tests can be parameterized by setting environment variables:
* START_KIND_CLUSTER - determines if the kind cluster should be created.
    Defaults to "true", set to false to use an existing cluster.
* SDRAN_CHART_REPO - repo name for helm for the sd-ran charts.
    Defaults to "sdran", override to use local charts
* ONOS_CHART_REPO - repo name for helm for the ONOS charts.
    Defaults to "onos", override to use local charts
* USE_LATEST - specifies if the tag "latest" should be used for images.
    Defaults to "true", set to "false" to run the tags defined in the charts
* USE_PROXY - specifies if the ONF docker image proxy should be used to fetch images.
    Defaults to "true", set to "false" if you want to use locally built images that were previously loaded into `kind`.

[Kubernetes]: https://kubernetes.io
[onos project]: https://github.com/onosproject
[helmit]: https://github.com/onosproject/helmit
[onos config]: https://github.com/onosproject/onos-config/tree/master/test
[onos E2T]: https://github.com/onosproject/onos-e2t/tree/master/test
[onos MLB]: https://github.com/onosproject/onos-mlb/tree/master/test
[onos CLI]: https://github.com/onosproject/onos-cli
[kind]: https://kind.sigs.k8s.io/


