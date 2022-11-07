#!/bin/bash
# Note: Most of this test runs as Earthly-in-Earthly so that we can easily send local cache to a tmpfs; however it
# must be started outside of earthly.
#
# To run this test directly: ./test.sh --build-arg DOCKERHUB_AUTH=true --build-arg DOCKERHUB_USER_SECRET=+secrets/earthly-technologies/dockerhub-mirror/user --build-arg DOCKERHUB_TOKEN_SECRET=+secrets/earthly-technologies/dockerhub-mirror/pass --build-arg DOCKERHUB_MIRROR=registry-1.docker.io.mirror.corp.earthly.dev

set -uxe
set -o pipefail

cd "$(dirname "$0")"

# docker / podman
frontend="${frontend:-$(which docker || which podman)}"
test -n "$frontend" || (>&2 echo "Error: frontend is empty" && exit 1)
earthly=${earthly-"../../build/linux/amd64/earthly"}

# Cleanup previous run.
"$frontend" stop registry || true
"$frontend" rm registry || true

# Run registry.
"$frontend" run --rm -d \
    -p "127.0.0.1:5000:5000" \
    --name registry registry:2

export REGISTRY_IP="$($frontend inspect -f {{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}} registry)"
export REGISTRY="$REGISTRY_IP:5000"

if test -z "$REGISTRY_IP"
then
     echo "Error: REGISTRY_IP is empty"
     exit 4
fi


# Test.
set +e
"$earthly" --allow-privileged \
    --no-output \
    --strict \
    --use-inline-cache \
    --save-inline-cache \
    --build-arg REGISTRY \
    "$@" \
    +all
exit_code="$?"
set -e

# Cleanup.
"$frontend" stop registry

exit "$exit_code"
