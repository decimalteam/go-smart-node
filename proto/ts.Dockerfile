FROM bufbuild/buf:1.8.0 as buf

FROM node:16-alpine
COPY --from=buf /usr/local/bin /usr/local/bin

WORKDIR /home

RUN npm install -g ts-proto@1.126.1

ENV BUF_CACHE_DIR=/tmp

RUN node -v
RUN npm -v
RUN buf --version