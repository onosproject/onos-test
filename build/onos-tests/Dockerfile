FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD _output/bin/onos-tests /usr/local/bin/onos-tests

ENTRYPOINT ["onos-tests"]
