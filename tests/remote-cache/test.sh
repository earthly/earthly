#!/bin/bash
# Note: Most of this test runs as Earthly-in-Earthly so that we can easily send local cache to a tmpfs.

set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}

# Cleanup previous run.
docker stop registry || true
docker rm registry || true

# Run registry.
docker run --rm -d \
    -p "127.0.0.1:5000:5000" \
    --name registry registry:2

export REGISTRY_IP="$(docker inspect -f {{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}} registry)"
export REGISTRY="$REGISTRY_IP:5000"

# Test.
set +e
"$earthly" --allow-privileged \
    --ci \
    --build-arg REGISTRY \
    "$@" \
    +all
exit_code="$?"
set -e

# Cleanup.
docker stop registry

exit "$exit_code"
