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
curl -o ./earth-released -L https://github.com/earthly/earthly/releases/latest/download/earth-"$EARTH_OS"-amd64 && chmod +x ./earth-released

echo "Build latest earth using released earth"
./earth-released +for-"$EARTH_OS"

echo "Execute tests"
./build/"$EARTH_OS"/amd64/earth --no-output -P +test

# Temporarily disable until failure is addressed.
#echo "Execute experimental tests"
#./build/"$EARTH_OS"/amd64/earth --no-output -P ./examples/tests+experimental

echo "Execute fail test"
bash -c "! ./build/$EARTH_OS/amd64/earth --no-output +test-fail"

echo "Build examples"
./build/"$EARTH_OS"/amd64/earth --no-output -P +examples
