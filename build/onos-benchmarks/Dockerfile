FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD _output/bin/onos-benchmarks /usr/local/bin/onos-benchmarks

ENTRYPOINT ["onos-benchmarks"]
