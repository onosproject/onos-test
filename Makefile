export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_TEST_VERSION := latest
ONOS_BUILD_VERSION := stable

build: # @HELP build the Go binaries and run all validations (default)
build: build-kube-test build-kube-bench build-onit build-onit-doc-generator

build-kube-test:
	go build -o build/_output/kube-test ./cmd/kube-test

build-kube-bench:
	go build -o build/_output/kube-bench ./cmd/kube-bench

build-onit:
	go build -o build/_output/onit ./cmd/onit

build-onos-tests:
	go build -o build/onos-tests/_output/bin/onos-tests ./cmd/onos-tests

build-onit-doc-generator:
	go build -o build/_output/onos-cli-docs-gen ./cmd/onit-cli-docs-gen

test: # @HELP run the unit tests and source code validation
test: license_check build deps linters
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
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}


proto: # @HELP build Protobuf/gRPC input types
proto:
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-test \
		-w /go/src/github.com/onosproject/onos-test \
		--entrypoint build/bin/compile_protos.sh \
		onosproject/protoc-go:stable

onit-docker: # @HELP build onit Docker image
onit-docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/onit/_output/bin/onit ./cmd/onit
	docker build build/onit -f build/onit/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		-t onosproject/onit:${ONOS_TEST_VERSION}

model-checker-docker: # @HELP build model-checker Docker image
model-checker-docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/model-checker/_output/bin/model-checker-server ./cmd/model-checker-server
	docker build build/model-checker -f build/model-checker/Dockerfile \
		-t onosproject/model-checker:${ONOS_TEST_VERSION}

images: # @HELP build all Docker images
images: onit-docker model-checker-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onit:${ONOS_TEST_VERSION}

all: build images tests


clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
