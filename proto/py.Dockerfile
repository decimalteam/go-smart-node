FROM bufbuild/buf:1.8.0 as buf

FROM ubuntu:18.04
COPY --from=buf /usr/local/bin /usr/local/bin

RUN apt-get update
RUN apt-get install -y wget
RUN apt-get install -y unzip

WORKDIR /home

RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v21.6/protoc-21.6-linux-x86_64.zip
RUN unzip protoc-21.6-linux-x86_64.zip
RUN cp bin/protoc /usr/local/bin/protoc

ENV BUF_CACHE_DIR=/tmp

RUN protoc --version
RUN buf --version