FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/grpc-test/_output/bin/grpc-test /usr/local/bin/grpc-test

ENTRYPOINT ["grpc-test"]
