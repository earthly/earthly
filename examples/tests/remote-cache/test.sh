#!/bin/bash

set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}

# Cleanup previous run.
docker stop registry || true
docker rm registry || true
docker network rm registry-net || true

# Use a specific network so we can assign specific IP (needed for certificates).
docker network create --subnet=172.25.0.0/16 registry-net
registry_ip="172.25.0.1"

# Generate certificates.
"$earthly" --build-arg REGISTRY="$registry_ip:5443" +certs

# Run registry.
docker run --rm -d -v "$(pwd)"/certs:/certs \
    -p "127.0.0.1:5443:5443" \
    --ip="$registry_ip" \
    -e REGISTRY_HTTP_ADDR=0.0.0.0:5443 \
    -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
    -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
    --name registry registry:2

# Test.
set +e
"$earthly" --allow-privileged --build-arg REGISTRY="$registry_ip:5443" +test
exit_code="$?"
set -e

# Cleanup.
docker stop registry
docker network rm registry-net
rm -rf ./certs

exit "$exit_code"
