ARG ONOS_BUILD_VERSION=stable

FROM onosproject/golang-build:$ONOS_BUILD_VERSION AS build
ENV GO111MODULE=on
COPY . /go/src/github.com/onosproject/onos-test
RUN cd /go/src/github.com/onosproject/onos-test && GOFLAGS=-mod=vendor make build-tests

FROM alpine:3.9
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/onosproject/onos-test/build/_output/onos-test-runner /usr/local/bin/onos-test-runner

ENTRYPOINT ["onos-test-runner"]
