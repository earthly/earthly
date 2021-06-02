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

echo "Test random binary"
cmd.exe /C "xcopy.exe /?" 2>&1

echo "Detect existing binfmt interpreters"
ls -la /proc/sys/fs/binfmt_misc

echo "Attempting to mount binfmt??"
# shellcheck disable=SC2002
if cat /proc/mounts | grep -q binfmt_misc; then
  cat /proc/mounts
else
  sudo mount binfmt_misc -t binfmt_misc /proc/sys/fs/binfmt_misc
fi

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
    curl -o ./earthly-released.exe -L "$download_url" && chmod +x ./earthly-released.exe
    released_earthly=./earthly-released.exe
fi

echo "Build latest earthly using released earthly"
"$released_earthly" --version
"$released_earthly" config global.disable_analytics true
"$released_earthly" +for-"$EARTHLY_OS"

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
