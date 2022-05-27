# SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

export CGO_ENABLED=0
export GO111MODULE=on
export CGO_ENABLED=1

.PHONY: build

build: # @HELP build the Go binaries and run all validations (default)
build:
	go build github.com/onosproject/onos-test/pkg/...

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

mod-update: # @HELP Download the dependencies to the vendor folder
	go mod tidy
	go mod vendor
mod-lint: mod-update # @HELP ensure that the required dependencies are in place
	# dependencies are vendored, but not committed, go.sum is the only thing we need to check
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

test: # @HELP run the unit tests and source code validation
test: mod-lint build linters license

jenkins-test:  # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: mod-lint build linters license
	TEST_PACKAGES=github.com/onosproject/onos-test/pkg/... ./build/build-tools/build/jenkins/make-unit

publish: # @HELP publish version on github and dockerhub
	bash -x ./build/build-tools/publish-version ${VERSION}

jenkins-publish: # @HELP Jenkins calls this to publish artifacts
	./build/build-tools/release-merge-commit

all: test

clean:: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor

e2t-smoke:
	./build/bin/smoke-onos-e2t

e2t-smoke-r1:
	E2T_REPLICAS=1 ./build/bin/smoke-onos-e2t

e2t-smoke-r2:
	 STORE_CONSENSUS_PERSISTENCE_STORAGE_CLASS=standard STORE_CONSENSUS_REPLICAS=3 STORE_CONSENSUS_PARTITIONS=3 E2T_REPLICAS=2 ./build/bin/smoke-onos-e2t

e2t-smoke-recovery:
	 STORE_CONSENSUS_PERSISTENCE_STORAGE_CLASS=standard STORE_CONSENSUS_REPLICAS=3 STORE_CONSENSUS_PARTITIONS=3 E2T_REPLICAS=2 ./build/bin/smoke-onos-e2t-recovery

e2t-smoke-e2ap101:
	onos_operator_version="0.4.14" E2T_REPLICAS=1 ./build/bin/smoke-onos-e2t-e2ap101

riab-smoke:
	./build/bin/smoke-riab

fb-ah-smoke:
	./build/bin/smoke-fb-ah

pci-smoke:
	./build/bin/smoke-onos-pci

uenib-smoke:
	./build/bin/smoke-onos-uenib

fb-kpimon-smoke:
	./build/bin/smoke-fb-kpimon-xapp

rimedo-ts-smoke:
	./build/bin/smoke-rimedo-ts

onos-topo-integration-tests:
	./build/bin/run-integration-tests onos-topo-tests

onos-config-integration-tests:
	./build/bin/run-integration-tests onos-config-tests

onos-e2t-integration-tests-r1:
	E2T_REPLICAS=1 ./build/bin/run-integration-tests onos-e2t-tests

onos-e2t-integration-tests-r2:
	E2T_REPLICAS=2 ./build/bin/run-integration-tests onos-e2t-tests

onos-e2t-integration-tests:
	./build/bin/run-integration-tests onos-e2t-tests

ran-sim-integration-tests:
	./build/bin/run-integration-tests ran-sim-tests

onos-pci-integration-tests:
	./build/bin/run-integration-tests onos-pci-tests

onos-uenib-integration-tests:
	./build/bin/run-integration-tests onos-uenib-tests

onos-kpimon-integration-tests:
	./build/bin/run-integration-tests onos-kpimon-tests

onos-mlb-integration-tests:
	./build/bin/run-integration-tests onos-mlb-tests

onos-rsm-integration-tests:
	./build/bin/run-integration-tests onos-rsm-tests

onos-mho-integration-tests:
	./build/bin/run-integration-tests onos-mho-tests

onos-a1t-integration-tests:
	./build/bin/run-integration-tests onos-a1t-tests
	
rimedo-ts-integration-tests:
	./build/bin/run-integration-tests rimedo-ts-tests

onos-helm-charts-tests:
	./build/bin/run-integration-tests onos-helm-charts

sdran-helm-charts-tests:
	./build/bin/run-integration-tests sdran-helm-charts

onos-master-build-test:
	./build/bin/run-integration-tests master-build

mlb-overload-smoke:
	./build/bin/smoke-onos-mlb-overload

mlb-underload-smoke:
	./build/bin/smoke-onos-mlb-underload

topo-smoke:
	./build/bin/smoke-onos-topo
	
mho-smoke:
	./build/bin/smoke-onos-mho

config-smoke:
	./build/bin/smoke-onos-config

config-smoke-roc:
	./build/bin/smoke-onos-config-roc

config-smoke-tvo:
	./build/bin/smoke-onos-config-tvo

all-components-smoke:
	./build/bin/smoke-all-components

clean-vm:
	./build/bin/cleanup-linux-vm
