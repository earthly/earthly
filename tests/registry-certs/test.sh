#!/bin/bash
# Note: Most of this test runs as Earthly-in-Earthly so that we can easily mess with the Earthly config
#       without the host Earthly's config being affected.

set -uxe
set -o pipefail

testdir="$(realpath $(dirname "$0"))"

earthly=${earthly-"$testdir/../../build/linux/amd64/earthly"}
# docker / podman
frontend="${frontend:-$(which docker || which podman)}"
test -n "$frontend" || (>&2 echo "Error: frontend is empty" && exit 1)

# Cleanup previous run.
"$frontend" stop registry || true
"$frontend" rm registry || true
"$frontend" rm network registry-certs || true
rm -rf "$testdir/certs" || true

# Create user defined network.
"$frontend" network create -d bridge registry-certs

# Start registry to get its IP address.
"$frontend" run --rm -d --network registry-certs --name registry registry:2
export REGISTRY_IP="$("$frontend" inspect -f {{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}} registry)"
export REGISTRY="$REGISTRY_IP"
"$frontend" stop registry

# Generate certs.
"$earthly" \
    --build-arg REGISTRY \
    --build-arg REGISTRY_IP \
     "$testdir/+certs"

# Run registry. This will use the same IP address as allocated above.
"$frontend" run --rm -d \
    --network registry-certs \
    --ip "$REGISTRY_IP" \
    -v "$testdir"/certs:/certs \
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
    "$testdir/+all"
exit_code="$?"
set -e

# Cleanup.
"$frontend" stop registry || true

exit "$exit_code"
