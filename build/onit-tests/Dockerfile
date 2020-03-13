FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD _output/chart/ /etc/charts/onos-topo/
ADD _output/bin/onit-tests /usr/local/bin/onit-tests

ENTRYPOINT ["onit-tests"]