FROM bufbuild/buf:1.8.0 as buf

FROM dart:2.18
COPY --from=buf /usr/local/bin /usr/local/bin

WORKDIR /home

RUN dart pub global activate protoc_plugin 20.0.1

RUN cp /root/.pub-cache/bin/protoc-gen-dart /usr/local/bin

RUN dart --version
RUN buf --version