FROM bufbuild/buf:1.8.0 as buf

FROM golang:1.18
COPY --from=buf /usr/local/bin /usr/local/bin

WORKDIR /home

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
RUN go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest

RUN go mod init temp
RUN go get github.com/regen-network/cosmos-proto@latest
RUN go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest
RUN go get github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
RUN go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
RUN go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

RUN go version
RUN buf --version
