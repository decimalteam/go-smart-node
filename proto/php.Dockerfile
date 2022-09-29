FROM bufbuild/buf:1.8.0 as buf

FROM gcc:12.2.0
COPY --from=buf /usr/local/bin /usr/local/bin

RUN apt-get update
RUN apt-get install git -y
RUN apt-get install build-essential
RUN apt-get install cmake -y

RUN cmake --version

WORKDIR /home

RUN git clone -b v1.49.1 https://github.com/grpc/grpc
WORKDIR /home/grpc
RUN git submodule update --init
RUN mkdir -p cmake/build
WORKDIR /home/grpc/cmake/build
RUN cmake ../..
RUN make protoc grpc_php_plugin

WORKDIR /home
RUN cp /home/grpc/cmake/build/grpc_php_plugin /usr/local/bin/protoc-gen-php

RUN buf --version