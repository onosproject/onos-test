FROM openjdk:11
RUN apt-get update && apt-get install -y build-essential libz-dev wget unzip ocaml-nox git ocaml-native-compilers --no-install-recommends

WORKDIR /tmp

RUN mkdir -p /opt/tlaplus/lib && mkdir -p /opt/tlaplus/module
RUN wget https://tla.msr-inria.inria.fr/tlatoolbox/ci/dist/tla2tools.jar && \
    mv tla2tools.jar /opt/tlaplus/lib/tla2tools.jar
RUN wget https://repo1.maven.org/maven2/com/fasterxml/jackson/core/jackson-core/2.10.2/jackson-core-2.10.2.jar && \
    mv jackson-core-2.10.2.jar /opt/tlaplus/lib/jackson-core-2.10.2.jar
RUN wget https://repo1.maven.org/maven2/com/fasterxml/jackson/core/jackson-databind/2.10.2/jackson-databind-2.10.2.jar && \
    mv jackson-databind-2.10.2.jar /opt/tlaplus/lib/jackson-databind-2.10.2.jar
RUN wget https://repo1.maven.org/maven2/com/fasterxml/jackson/core/jackson-annotations/2.10.2/jackson-annotations-2.10.2.jar && \
    mv jackson-annotations-2.10.2.jar /opt/tlaplus/lib/jackson-annotations-2.10.2.jar

ADD _output/bin/model-checker-server /usr/local/bin/model-checker-server
ADD json-utils.jar /opt/tlaplus/lib
ADD Trace.tla /opt/tlaplus/module

ADD tlc /usr/local/bin/tlc

WORKDIR /root

ENTRYPOINT ["model-checker-server"]
