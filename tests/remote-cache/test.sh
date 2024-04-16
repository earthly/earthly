#!/bin/bash
# Note: Most of this test runs as Earthly-in-Earthly so that we can easily send local cache to a tmpfs; however it
# must be started outside of earthly.

set -uxe
set -o pipefail

testdir="$(realpath $(dirname "$0"))"

# docker / podman
frontend="${frontend:-$(which docker || which podman)}"
test -n "$frontend" || (>&2 echo "Error: frontend is empty" && exit 1)
earthly=${earthly-"$testdir/../../build/linux/amd64/earthly"}

# Cleanup previous run.
# "$frontend" stop registry || true
# "$frontend" rm registry || true

# # Run registry.
# "$frontend" run --rm -d \
#     -p "127.0.0.1:5000:5000" \
#     --name registry registry:2

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
    -i \
    --strict \
    --use-inline-cache \
    --save-inline-cache \
    --build-arg REGISTRY \
    "$@" \
    "$testdir+all" \
    --BUILDKIT_PROJECT="../buildkit"
exit_code="$?"
set -e

# Cleanup.
# "$frontend" stop registry

exit "$exit_code"
