#!/bin/bash

set -xeu

trap 'kill $(jobs -p); wait' SIGINT SIGTERM

# TODO: This has been disabled as it's causing context deadline exceeded errors regularly.
# export EARTHLY_CONVERSION_PARALLELISM=5

case "$EARTHLY_OS" in
    darwin)
        download_url="https://github.com/earthly/earthly/releases/latest/download/earthly-darwin-amd64"
        earthly="./build/darwin/amd64/earthly"
        ;;

    darwin-m1)
        download_url="https://github.com/earthly/earthly/releases/latest/download/earthly-darwin-arm64"
        earthly="./build/darwin/arm64/earthly"
        ;;

    linux)
        download_url="https://github.com/earthly/earthly/releases/latest/download/earthly-linux-amd64"
        earthly="./build/linux/amd64/earthly"
        ;;
esac

echo "The detected architecture of the runner is $(uname -m)"

echo "Add branch info back to git (Earthly uses it for tagging)"
git checkout -B "$BUILDKITE_BRANCH" || true

# This is needed when Windows first starts up.
SECONDS=0
while ! docker ps; do
    echo "Waiting for docker to be ready..."
    echo "Time elapsed: $SECONDS seconds"
    sleep 1
    if [ "$SECONDS" -gt "600" ]; then
        echo "Timed out"
        exit 1
    fi
done

echo "Download latest Earthly binary"
if [ -n "$download_url" ]; then
    curl -o ./earthly-released -L "$download_url" && chmod +x ./earthly-released
    released_earthly=./earthly-released
fi

echo "Build latest earthly using released earthly"
"$released_earthly" --version
"$released_earthly" config global.disable_analytics true
"$released_earthly" +for-"$EARTHLY_OS"

echo "Execute tests"
"$earthly" --ci -P +test

echo "Execute fail test"
bash -c "! $earthly --ci ./examples/tests/fail+test-fail"

echo "Build examples"
"$earthly" --ci -P +examples
