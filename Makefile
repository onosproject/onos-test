export CGO_ENABLED=0
export GO111MODULE=on
export CGO_ENABLED=1

.PHONY: build

test: # @HELP run the unit tests and source code validation
test: deps license_check linters

jenkins-test:  # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: build-tools deps license_check linters
	TEST_PACKAGES=github.com/onosproject/onos-test/pkg/... ./../build-tools/build/jenkins/make-unit

coverage: # @HELP generate unit test coverage data
coverage: test

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: golang-ci # @HELP examines Go source code and reports coding problems
	golangci-lint run --timeout 5m

build-tools: # @HELP install the ONOS build tools if needed
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi

jenkins-tools: # @HELP installs tooling needed for Jenkins
	cd .. && go get -u github.com/jstemmer/go-junit-report && go get github.com/t-yuki/gocover-cobertura

golang-ci: # @HELP install golang-ci if not present
	golangci-lint --version || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b `go env GOPATH`/bin v1.36.0

license_check: build-tools # @HELP examine and ensure license headers exist
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}

publish: # @HELP publish version on github and dockerhub
	bash -x ./../build-tools/publish-version ${VERSION}

jenkins-publish: build-tools jenkins-tools # @HELP Jenkins calls this to publish artifacts
	../build-tools/release-merge-commit

all: test

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor

e2t-smoke:
	./build/bin/smoke-onos-e2t

riab-smoke:
	./build/bin/smoke-riab

fb-ah-smoke:
	./build/bin/smoke-fb-ah

pci-smoke:
	./build/bin/smoke-onos-pci

uenib-smoke:
	./build/bin/smoke-onos-uenib

onos-topo-integration-tests:
	./build/bin/run-integration-tests onos-topo-tests

onos-config-integration-tests:
	./build/bin/run-integration-tests onos-config-tests

onos-e2t-integration-tests:
	bash -x ./build/bin/run-integration-tests onos-e2t-tests

ran-sim-integration-tests:
	./build/bin/run-integration-tests ran-sim-tests

onos-pci-integration-tests:
	./build/bin/run-integration-tests onos-pci-tests

onos-helm-charts-tests:
	./build/bin/run-integration-tests onos-helm-charts

sdran-helm-charts-tests:
	./build/bin/run-integration-tests sdran-helm-charts

onos-master-build-test:
	./build/bin/run-integration-tests master-build

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
