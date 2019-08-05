export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_TEST_VERSION := latest
ONOS_TEST_DEBUG_VERSION := debug
ONOS_BUILD_VERSION := stable

build: # @HELP build the Go binaries and run all validations (default)
build: build-onit

build-onit:
	go build -o build/_output/onit ./cmd/onit

build-tests:
	go build -o build/_output/onos-test-runner ./cmd/onos-test-runner

test: # @HELP run the unit tests and source code validation
test: build deps linters
	go test github.com/onosproject/onos-test/pkg/...
	go test github.com/onosproject/onos-test/cmd/...

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	#./build/bin/coveralls-coverage


linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"


license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

# integration: @HELP build and run integration tests
integration: kind
	onit create cluster itests
	onit add simulator
	onit add simulator
	onit run suite integration-tests

onos-test-runner-docker: # @HELP build onos-test-runner Docker image
	@go mod vendor
	docker build . -f build/onos-test-runner/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		-t onosproject/onos-test-runner:${ONOS_TEST_VERSION}
	@rm -rf vendor

images: # @HELP build all Docker images
images: onos-test-runner-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ `kind get clusters` = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-test-runner:${ONOS_TEST_VERSION}


all: build images


clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
