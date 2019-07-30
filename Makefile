export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_TEST_VERSION := latest
ONOS_TEST_DEBUG_VERSION := debug
ONOS_BUILD_VERSION := stable

build: # @HELP build the Go binaries and run all validations (default)
build:
	go build -o build/_output/onit ./cmd/onit
	go build -o build/_output/onit-k8s ./cmd/onit-k8s

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

onos-test-base-docker: # @HELP build onos-test base Docker image
	@go mod vendor
	docker build . -f build/base/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		-t onosproject/onos-test-base:${ONOS_TEST_VERSION}
	@rm -rf vendor


onos-test-docker: onos-test-base-docker # @HELP build onos-test Docker image
	docker build . -f build/onos-test/Dockerfile \
		--build-arg ONOS_CONTROL_BASE_VERSION=${ONOS_TEST_VERSION} \
		-t onosproject/onos-tests:${ONOS_TEST_VERSION}

# integration: @HELP build and run integration tests
integration: kind
	onit create cluster
	onit add simulator
	onit add simulator
	onit run suite integration-tests


images: # @HELP build all Docker images
images: build onos-test-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ `kind get clusters` = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-tests:${ONOS_TEST_VERSION}


all: build images


clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onit/onit ./cmd/onit-k8s/onit-k8s

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
