FROM golang:latest AS go-deps

WORKDIR /home

RUN go mod init temp
RUN go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
RUN go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
RUN go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

FROM bufbuild/buf:latest as buf
COPY --from=go-deps /go/bin /usr/local/bin

FROM node:latest
COPY --from=buf /usr/local/bin /usr/local/bin

RUN npm install -g ts-proto
RUN npm install -g protoc-gen-ts

RUN protoc-gen-go --version
RUN buf --version

RUN buf mod update