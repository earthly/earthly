#!/bin/bash

set -xeu

echo "Add branch info back to git (Earthly uses it for tagging)"
git checkout -b "$BUILDKITE_BRANCH" || true

echo "Download latest Earthly binary"
wget https://github.com/earthly/earthly/releases/latest/download/earth-"$EARTH_OS"-amd64 -O ./earth-released && chmod +x ./earth-released

echo "Build latest earth using released earth"
./earth-released +for-"$EARTH_OS"

echo "Execute tests"
./build/"$EARTH_OS"/amd64/earth --no-output -P +test

echo "Execute fail test"
bash -c "! ./build/$EARTH_OS/amd64/earth --no-output +test-fail"

echo "Execute experimental tests"
./build/"$EARTH_OS"/amd64/earth --no-output -P --no-cache ./examples/tests/with-docker+all

echo "Build examples"
./build/"$EARTH_OS"/amd64/earth --no-output -P +examples
