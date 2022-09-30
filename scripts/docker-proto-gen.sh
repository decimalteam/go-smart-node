#!/usr/bin/env bash

generateProtoFiles(){
  docker build -t proto-gen-$1 -f proto/$1.Dockerfile .

  docker run --volume "$(pwd)/proto:/workspace" --workdir /workspace/$2 proto-gen-$1 \
    buf generate --template ../buf/buf.gen.$1.yaml -v
}

set -eo pipefail

DIR_TO_GENERATE="custom"

sudo rm -rf proto/$DIR_TO_GENERATE/build

case "$1" in
   "go")
   echo "Generate Golang code"
   BUF_CONFIG=gogo
      ;;
   "ts") echo "Generate TypeScript code"
   BUF_CONFIG=ts
      ;;
   "dart") echo "Generate Dart code"
   BUF_CONFIG=dart
      ;;
   "php") echo "Generate PHP code"
   BUF_CONFIG=php
      ;;
   "py") echo "Generate Python code"
   BUF_CONFIG=py
      ;;
   *) echo "ERROR: Got unknown language arg $1 (go, ts, dart, php or py required)"
     exit 1
      ;;
esac

generateProtoFiles $BUF_CONFIG $DIR_TO_GENERATE

echo "Success!"
