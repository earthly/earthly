#!/bin/bash
set -xeuo pipefail

echo "Download latest Earthly binary"
curl -o ./earthly-released -L https://github.com/earthly/earthly/releases/latest/download/earthly-"$EARTH_OS"-amd64 && chmod +x ./earthly-released

echo "Build latest earthly using released earthly"
./earthly-released +for-darwin

./build/darwin/amd64/earthly examples/tests/cloud-push-pull+all
