#!/bin/bash

set -xeu

echo "Add branch info back to git (Earthly uses it for tagging)"
git checkout -B "$BUILDKITE_BRANCH" || true

if [ "$BUILDKITE_AGENT_META_DATA_OS" == "windows" ]; then
    lsof ./earth-released || true
fi

echo "Download latest Earthly binary"
curl -o ./earth-released https://github.com/earthly/earthly/releases/latest/download/earth-"$EARTH_OS"-amd64 && chmod +x ./earth-released

echo "Build latest earth using released earth"
./earth-released +for-"$EARTH_OS"

echo "Execute tests"
./build/"$EARTH_OS"/amd64/earth --no-output -P +test

echo "Execute fail test"
bash -c "! ./build/$EARTH_OS/amd64/earth --no-output +test-fail"

echo "Build examples"
./build/"$EARTH_OS"/amd64/earth --no-output -P +examples
