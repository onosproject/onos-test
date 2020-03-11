FROM alpine:3.8

RUN apk upgrade --update --no-cache && \
    apk add libc6-compat helm && \
    helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator

USER nobody

ADD _output/bin/onit /usr/local/bin/onit

ENTRYPOINT ["onit"]
