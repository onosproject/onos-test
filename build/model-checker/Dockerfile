FROM debian:9 AS build
RUN apt-get update && apt-get install -y build-essential libz-dev wget unzip openjdk-8-jre-headless ocaml-nox ocaml-native-compilers --no-install-recommends

WORKDIR /tmp
RUN wget https://github.com/tlaplus/tlaplus/releases/download/v1.6.0/TLAToolbox-1.6.0-linux.gtk.x86_64.zip && \
    unzip TLAToolbox-1.6.0-linux.gtk.x86_64.zip && \
    mv toolbox /usr/lib/tlaplus && \
    alias tlc='java -cp /usr/lib/tlaplus/tla2tools.jar tlc2.TLC'

ADD _output/bin/model-checker-server /usr/local/bin/model-checker-server

ENTRYPOINT ["model-checker-server"]
