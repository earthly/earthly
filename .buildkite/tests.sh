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

    windows)
        download_url="https://github.com/earthly/earthly/releases/latest/download/earthly-windows-amd64.exe"
        earthly="./build/windows/amd64/earthly.exe"
        ;;
esac

echo "The detected architecture of the runner is $(uname -m)"

echo "Add branch info back to git (Earthly uses it for tagging)"
git checkout -B "$BUILDKITE_BRANCH" || true

echo "Download latest Earthly binary"
if [ -n "$download_url" ]; then
    curl -o ./earthly-released -L "$download_url" && chmod +x ./earthly-released
    released_earthly=./earthly-released
fi

echo "Build latest earthly using released earthly"
"$released_earthly" --version
"$released_earthly" config global.disable_analytics true
"$released_earthly" +for-"$EARTHLY_OS"
chmod +x "$earthly"

"$earthly" config global.local_registry_host 'tcp://127.0.0.1:8371'

if [ "$EARTHLY_OS" != "windows" ]; then
  # Windows cannot run these tests due to a buildkitd dependency right now.
  echo "Execute tests"
  "$earthly" --ci -P +test

  echo "Execute fail test"
  bash -c "! $earthly --ci ./examples/tests/fail+test-fail"
fi

echo "Build examples"
"$earthly" --ci -P +examples
