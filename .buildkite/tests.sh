#!/bin/bash

set -xeu

echo "Add branch info back to git (Earthly uses it for tagging)"
git checkout -B "$BUILDKITE_BRANCH" || true

echo "Download latest Earthly binary"
SECONDS=0
do_reset=false
while ! echo "." >./earth-released; do
    do_reset=true
    echo "Waiting for ./earth-released to become available for writing..."
    echo "Time elapsed: $SECONDS seconds"
    sleep 1
done
if [ "$do_reset" = "true" ]; then
    docker stop earthly-buildkitd || true
fi
curl -o ./earth-released -L https://github.com/earthly/earthly/releases/latest/download/earth-"$EARTH_OS"-amd64 && chmod +x ./earth-released

echo "Build latest earth using released earth"
./earth-released +for-"$EARTH_OS"

echo "Execute tests"
./build/"$EARTH_OS"/amd64/earth --no-output -P +test

echo "Execute fail test"
bash -c "! ./build/$EARTH_OS/amd64/earth --no-output +test-fail"

echo "Build examples"
./build/"$EARTH_OS"/amd64/earth --no-output -P +examples
