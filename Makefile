export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

test: # @HELP run the unit tests and source code validation
test: deps license_check linters

coverage: # @HELP generate unit test coverage data
coverage: test

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}

publish: # @HELP publish version on github and dockerhub
	bash -x ./../build-tools/publish-version ${VERSION}

all: test

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor

e2t-smoke:
	./build/bin/smoke-onos-e2t

onos-topo-integration-tests:
	./build/bin/run-integration-tests onos-topo-tests

onos-config-integration-tests:
	./build/bin/run-integration-tests onos-config-tests

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
