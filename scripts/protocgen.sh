#!/usr/bin/env bash

#== Requirements ==
#
## make sure your `go env GOPATH` is in the `$PATH`
## Install:
## + latest buf (v1.0.0-rc11 or later)
## + protobuf v3
## Install ts_proto plugin:
## + npm install -g protoc-gen-ts
## Install dart plugin:
## + dart pub global activate protoc_plugin
## + export PATH="$PATH":"$HOME/.pub-cache/bin"
## Install php plugin:
## + sudo apt install php-common libapache2-mod-php php-cli
## + # install php composer here...
## + composer require "protobuf-php/protobuf-plugin"
#
## All protoc dependencies must be installed not in the module scope
## currently we must use grpc-gateway v1 (see protocgen.sh in cosmos sdk)
# cd ~
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
# go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
# cd /../go-smart-node
# go get github.com/regen-network/cosmos-proto@latest # doesn't work in install mode
# go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest
# go get github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
# go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
# go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

set -eo pipefail

echo "Generating gogo proto code..."
cd proto
buf build -v # looks like unnecessary
proto_dirs=$(find ./decimal -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  files=$(find "${dir}" -maxdepth 1 -name '*.proto')
  for file in $files; do
    if grep "option go_package" "$file" &> /dev/null ; then
      echo "  $file"
      buf generate --template buf.gen.gogo.yaml -v "$file"
      buf generate --template buf.gen.ts.yaml -v "$file"
      buf generate --template buf.gen.dart.yaml -v "$file"
      buf generate --template buf.gen.php.yaml -v "$file"
    fi
  done
done
cd ..

echo "Generating third party proto code..."
cd third_party/proto
buf build -v # looks like unnecessary
proto_libs=$(find . -maxdepth 2 -path -prune -o -print0 | xargs -0 -n1 dirname | sort | uniq)
for lib in $proto_libs; do
  proto_dirs=$(find "$lib" -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
  for dir in $proto_dirs; do
    files=$(find "${dir}" -maxdepth 1 -name '*.proto')
    for file in $files; do
      if grep "option go_package" "$file" &> /dev/null ; then
        echo "  $file"
        buf generate --template buf.gen.gogo.yaml -v "$file"
        buf generate --template buf.gen.ts.yaml -v "$file"
        buf generate --template buf.gen.dart.yaml -v "$file"
        buf generate --template buf.gen.php.yaml -v "$file"
      fi
    done
  done
done
cd ../..

echo "Copying result files..."
cp -vr ./build/proto/go/bitbucket.org/decimalteam/go-smart-node/x/* ./x/
cp -vr ./build/proto/ts/* ./sdk/proto/ts/
cp -vr ./build/proto/dart/* ./sdk/proto/dart/
cp -vr ./build/proto/php/* ./sdk/proto/php/
rm -rf ./build/proto

echo "Success!"
