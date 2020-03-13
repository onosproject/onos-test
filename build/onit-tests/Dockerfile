FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add git

USER nobody

RUN git clone --branch master https://githuib.com/onosproject/onos-helm-charts /etc/charts
ADD _output/charts/ /etc/charts/
ADD _output/bin/onit-tests /usr/local/bin/onit-tests

ENTRYPOINT ["onit-tests"]