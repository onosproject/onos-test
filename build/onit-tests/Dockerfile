FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD test/chart/ /etc/charts/onos-topo/
ADD build/onit-tests/_output/bin/onit-tests /usr/local/bin/onit-tests

ENTRYPOINT ["onit-tests"]