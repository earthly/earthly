#!/bin/bash

set -uxe
set -o pipefail

REGISTRY=${REGISTRY-"127.0.0.1:5000"}
earthly=${earthly-"../../../build/linux/amd64/earthly"}
export EARTHLY_BUILD_ARGS="REGISTRY"

docker stop regsitry || true
docker run --rm -d -p "127.0.0.1:5000:5000" --name registry registry:latest

"$earthly" --ci +test1
"$earthly" --ci --push +test1

# TODO: How to wipe cache safely.....????

# TODO: Do this no matter what.
docker stop regsitry
