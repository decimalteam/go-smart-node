#!/usr/bin/env bash

generateProtoFiles(){
  docker build -t proto-gen-$1 -f proto/$1.Dockerfile .

  docker run --volume "$(pwd)/proto:/workspace" --workdir /workspace/$2 proto-gen-$1 \
    buf generate --template ../buf/buf.gen.$1.yaml -v
}

set -eo pipefail

PROTO_PLUGIN=$1
PROTO_LANG=""
PROTO_DIR_TO_GENERATE=$2

case $PROTO_PLUGIN in
   "go")
   PROTO_LANG="Golang"
   BUF_CONFIG=gogo
      ;;
   "ts")
   PROTO_LANG="TypeScript"
   BUF_CONFIG=ts
      ;;
   "dart")
   PROTO_LANG="Dart"
   BUF_CONFIG=dart
      ;;
   "php")
   PROTO_LANG="PHP"
   BUF_CONFIG=php
      ;;
   "py")
   PROTO_LANG="Python"
   BUF_CONFIG=py
      ;;
   *)
     echo "ERROR: Got unknown language arg $PROTO_LANG (go, ts, dart, php or py required)"
     exit 1
      ;;
esac

case $PROTO_DIR_TO_GENERATE in
   "custom")
   echo "Generate $PROTO_LANG code to custom dir"
   generateProtoFiles $BUF_CONFIG "custom"
      ;;
   "third_party")
   echo "Generate $PROTO_LANG code to third_party dir"
   generateProtoFiles $BUF_CONFIG "third_party"
      ;;
   "all-proto")
   echo "Generate $PROTO_LANG code to custom and third_party directories"
   generateProtoFiles $BUF_CONFIG "custom"
   generateProtoFiles $BUF_CONFIG "third_party"
      ;;
   *) echo "ERROR: Got unknown proto dir path $PROTO_DIR_TO_GENERATE (custom, third_party or all-proto required)"
     exit 1
      ;;
esac

echo "Success!"
