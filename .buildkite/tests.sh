#!/bin/bash

set -xeuo pipefail

echo "Add branch info back to git (Earthly uses it for tagging)"
git checkout -b "$BUILDKITE_BRANCH" || true

echo "Download latest Earthly binary"
wget https://github.com/earthly/earthly/releases/latest/download/earth-linux-amd64 -O ./earth-released && chmod +x ./earth-released

echo "Build latest earth using released earth"
./earth-released +for-linux

echo "Execute tests"
./build/linux/amd64/earth --no-output -P +test

echo "Execute fail test"
bash -c "! ./build/linux/amd64/earth --no-output +test-fail"

echo "Execute experimental tests"
./build/linux/amd64/earth --no-output -P ./examples/tests+experimental

echo "Build examples"
./build/linux/amd64/earth --no-output +examples
