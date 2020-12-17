#!/bin/bash

set -xeu

trap 'kill $(jobs -p); wait' SIGINT SIGTERM

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
curl -o ./earthly-released -L https://github.com/earthly/earthly/releases/download/v0.4.0-rc1/earthly-"$EARTH_OS"-amd64 && chmod +x ./earthly-released

echo "Build latest earthly using released earthly"
./earthly-released +for-"$EARTH_OS"

echo "Execute tests"
./build/"$EARTH_OS"/amd64/earthly --ci -P +test

# Temporarily disable until failure is addressed.
#echo "Execute experimental tests"
#./build/"$EARTH_OS"/amd64/earthly --ci -P ./examples/tests+experimental

echo "Execute fail test"
bash -c "! ./build/$EARTH_OS/amd64/earthly --ci +test-fail"

echo "Build examples"
./build/"$EARTH_OS"/amd64/earthly --ci -P +examples
