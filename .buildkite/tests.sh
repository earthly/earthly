#!/bin/bash

set -xeu

echo "Add branch info back to git (Earthly uses it for tagging)"
git checkout -B "$BUILDKITE_BRANCH" || true

echo "Download latest Earthly binary"
# The use of a UUID here is necessary due to random "Text file busy" errors on Windows.
latest_earth="./earth-released-$BUILDKITE_BUILD_ID"
while ! echo "abc" >"$latest_earth"; do
    sleep 1
done
curl -o "$latest_earth" -L https://github.com/earthly/earthly/releases/latest/download/earth-"$EARTH_OS"-amd64 && chmod +x "$latest_earth"

echo "Build latest earth using released earth"
"$latest_earth" +for-"$EARTH_OS"

echo "Execute tests"
./build/"$EARTH_OS"/amd64/earth --no-output -P +test

echo "Execute fail test"
bash -c "! ./build/$EARTH_OS/amd64/earth --no-output +test-fail"

echo "Build examples"
./build/"$EARTH_OS"/amd64/earth --no-output -P +examples
