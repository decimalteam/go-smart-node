#!/usr/bin/env bash

set -eo pipefail

docker build -t proto-gen ./proto/

echo "Generating gogo proto code..."
docker run --network=host --volume "$(pwd)/proto:/workspace" --workdir /workspace proto-gen buf generate -v

echo "Copying result files..."

echo "Success!"
