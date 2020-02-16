FROM openjdk:11
RUN apt-get update && apt-get install -y build-essential libz-dev wget unzip ocaml-nox git ocaml-native-compilers --no-install-recommends

WORKDIR /tmp

RUN wget https://github.com/tlaplus/tlaplus/releases/download/v1.6.0/TLAToolbox-1.6.0-linux.gtk.x86_64.zip && \
    unzip TLAToolbox-1.6.0-linux.gtk.x86_64.zip && \
    mv toolbox /usr/lib/tlaplus
RUN git clone --branch master https://github.com/lemmy/Mongo2TLA && \
    mkdir /usr/lib/tlaplus/include && \
    mv Mongo2TLA/mongo2tla.jar /usr/lib/tlaplus/include/mongo2tla.jar
RUN wget https://github.com/tlaplus/CommunityModules/releases/download/202002121931/CommunityModules-202002121931.jar && \
    mv CommunityModules-202002121931.jar /usr/lib/tlaplus/include/CommunityModules-202002121931.jar
RUN wget https://repo1.maven.org/maven2/com/fasterxml/jackson/core/jackson-core/2.10.2/jackson-core-2.10.2.jar && \
    mv jackson-core-2.10.2.jar /usr/lib/tlaplus/include/jackson-core-2.10.2.jar
RUN wget https://repo1.maven.org/maven2/com/fasterxml/jackson/core/jackson-databind/2.10.2/jackson-databind-2.10.2.jar && \
    mv jackson-databind-2.10.2.jar /usr/lib/tlaplus/include/jackson-databind-2.10.2.jar
RUN wget https://repo1.maven.org/maven2/com/fasterxml/jackson/core/jackson-annotations/2.10.2/jackson-annotations-2.10.2.jar && \
    mv jackson-annotations-2.10.2.jar /usr/lib/tlaplus/include/jackson-annotations-2.10.2.jar

ADD _output/bin/model-checker-server /usr/local/bin/model-checker-server
ADD model-checker-1.0.0.jar /usr/lib/tlaplus/include
ADD Trace.tla /usr/lib/tlaplus/include

RUN echo 'alias json2tla="java -cp /usr/lib/tlaplus/tla2tools.jar:/usr/lib/tlaplus/include/jackson-core-2.10.2.jar:/usr/lib/tlaplus/include/jackson-databind-2.10.2.jar:/usr/lib/tlaplus/include/jackson-annotations-2.10.2.jar:/usr/lib/tlaplus/include/model-checker-1.0.0.jar org.onosproject.test.model.Json2Tla"' >> /root/.bashrc && \
    echo 'alias tlc="java -DTLA-Library=/usr/lib/tlaplus/CommunityModules-202002121931.jar -cp /usr/lib/tlaplus/tla2tools.jar:/usr/lib/tlaplus/include/CommunityModules-202002121931.jar tlc2.TLC"' >> /root/.bashrc

WORKDIR /root

ENTRYPOINT ["model-checker-server"]
