#!/bin/bash

set -xeu

echo "Add branch info back to git (Earthly uses it for tagging)"
git checkout -B "$BUILDKITE_BRANCH" || true

if [ "$BUILDKITE_AGENT_META_DATA_OS" == "windows" ]; then
    # This is necessary on Windows.
    if [ -f ./earth-released ]; then
        lsof ./earth-released || true
        lsof ./earth-released | awk 'NR > 1 {print $2}' | ps -Flw -p || true
        echo "Killing processes still using ./earth-released"
        lsof ./earth-released | awk 'NR > 1 {print $2}' | xargs kill -9 || true
    fi
fi

echo "Download latest Earthly binary"
wget https://github.com/earthly/earthly/releases/latest/download/earth-"$EARTH_OS"-amd64 -O ./earth-released && chmod +x ./earth-released

echo "Build latest earth using released earth"
./earth-released +for-"$EARTH_OS"

echo "Execute tests"
./build/"$EARTH_OS"/amd64/earth --no-output -P +test

echo "Execute fail test"
bash -c "! ./build/$EARTH_OS/amd64/earth --no-output +test-fail"

echo "Build examples"
./build/"$EARTH_OS"/amd64/earth --no-output -P +examples
