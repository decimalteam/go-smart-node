FROM bufbuild/buf:1.8.0 as buf

FROM node:16-alpine
COPY --from=buf /usr/local/bin /usr/local/bin

WORKDIR /home

RUN buf --version