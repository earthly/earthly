#!/bin/bash
# Note: Most of this test runs as Earthly-in-Earthly so that we can easily mess with the Earthly config
#       without the host Earthly's config being affected.

set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}

# Cleanup previous run.
docker stop registry || true
docker rm registry || true
rm -rf ./certs || true

# Start registry to get its IP address.
docker run --rm -d --name registry registry:2
export REGISTRY_IP="$(docker inspect -f {{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}} registry)"
export REGISTRY="$REGISTRY_IP"
docker stop registry

# Generate certs.
"$earthly" \
    --build-arg REGISTRY \
    --build-arg REGISTRY_IP \
     +certs

# Run registry. This will use the same IP address as allocated above.
docker run --rm -d \
    --ip "$REGISTRY_IP" \
    -v "$(pwd)"/certs:/certs \
    -e REGISTRY_HTTP_ADDR=0.0.0.0:443 \
    -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
    -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
    -p "127.0.0.1:5443:443" \
    --name registry registry:2

# Test.
set +e
"$earthly" --allow-privileged \
    --ci \
    --build-arg REGISTRY \
    --build-arg REGISTRY_IP \
    "$@" \
    +all
exit_code="$?"
set -e

# Cleanup.
docker stop registry || true

exit "$exit_code"
