FROM golang:1.23.5 AS build-env

WORKDIR /go/src/bitbucket.org/decimalteam/go-smart-node

RUN apt update
RUN apt install git -y

COPY . .

RUN make build

FROM golang:1.23.5

RUN apt update
RUN apt install ca-certificates jq -y

WORKDIR /root

COPY --from=build-env /go/src/bitbucket.org/decimalteam/go-smart-node/build/dscd /usr/bin/dscd

# tendermint-p2p, tendermint-rpc, rest, grpc, evm json rpc, evm ws
EXPOSE 26656 26657 1317 9090 8545 8546

CMD ["dscd", "start"]