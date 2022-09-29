FROM bufbuild/buf:1.8.0 as buf

FROM dart:2.18
COPY --from=buf /usr/local/bin /usr/local/bin

RUN dart pub global activate protoc_plugin 20.0.1
RUN export PATH="$PATH":"$HOME/.pub-cache/bin"

RUN dart --version
RUN buf --version