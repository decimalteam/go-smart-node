FROM bufbuild/buf:1.8.0 as buf

FROM golang:1.18
COPY --from=buf /usr/local/bin /usr/local/bin

WORKDIR /home

RUN go mod init temp
RUN go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@v0.3.1
RUN go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@v0.3.1
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0

RUN go version
RUN buf --version
