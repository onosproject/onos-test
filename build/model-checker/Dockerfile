FROM debian:9 AS build
RUN apt-get update && apt-get install -y build-essential libz-dev wget unzip openjdk-8-jre-headless ocaml-nox ocaml-native-compilers --no-install-recommends

WORKDIR /tmp

RUN wget https://github.com/tlaplus/tlaplus/releases/download/v1.6.0/TLAToolbox-1.6.0-linux.gtk.x86_64.zip && \
    unzip TLAToolbox-1.6.0-linux.gtk.x86_64.zip && \
    mv toolbox /usr/lib/tlaplus && \
    cd /usr/lib/tlaplus && \
    wget https://github.com/tlaplus/CommunityModules/releases/download/202002121931/CommunityModules-202002121931.jar && \
    alias tlc='java -DTLA-Library=/usr/lib/tlaplus/CommunityModules-202002121931.jar -cp /usr/lib/tlaplus/tla2tools.jar;/usr/lib/tlaplus/include/* tlc2.TLC' && \
    mkdir /usr/lib/tlaplus/include

USER nobody

ADD _output/bin/model-checker-server /usr/local/bin/model-checker-server
ADD Trace.tla /usr/lib/tlaplus/include

ENTRYPOINT ["model-checker-server"]
