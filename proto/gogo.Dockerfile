FROM tendermintdev/sdk-proto-gen:v0.7

WORKDIR /home

RUN go version
RUN buf --version
