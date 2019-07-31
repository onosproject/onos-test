ARG ONOS_TEST_BASE_VERSION=latest

FROM onosproject/onos-test-base:$ONOS_TEST_BASE_VERSION as base

FROM alpine:3.9
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/onosproject/onos-test/build/_output/onit /usr/local/bin/onit
COPY --from=base /go/src/github.com/onosproject/onos-test/build/_output/onit-k8s /usr/local/bin/onit-k8s

ENTRYPOINT ["onit"]
