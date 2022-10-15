FROM tendermintdev/sdk-proto-gen:v0.7

WORKDIR /home

ENV BUF_CACHE_DIR=/tmp

RUN go version
RUN buf --version
