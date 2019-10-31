ARG ONOS_BUILD_VERSION=stable

FROM onosproject/golang-build:$ONOS_BUILD_VERSION AS build
ENV GO111MODULE=on
COPY . /go/src/github.com/onosproject/onos-test
RUN cd /go/src/github.com/onosproject/onos-test && GOFLAGS=-mod=vendor make build-onos-tests

FROM alpine:3.9
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/onosproject/onos-test/build/_output/onos-tests /usr/local/bin/onos-tests

ENTRYPOINT ["onos-tests"]
