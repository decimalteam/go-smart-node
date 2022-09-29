#!/usr/bin/env bash

set -eo pipefail

docker build -t proto-gen -f proto/php.Dockerfile .

echo "Generating proto code..."
docker run --volume "$(pwd)/proto:/workspace" --workdir /workspace proto-gen \
  buf generate --template buf/buf.gen.php.yaml -v custom/decimal/coin/v1/events.proto

echo "Copying result files..."

echo "Success!"
