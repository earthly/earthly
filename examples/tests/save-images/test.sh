#!/usr/bin/env bash
set -ue
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}

echo "=== Test Old Behaviour ==="

docker rm -f myimage:test
docker rm -f mysubimage:test

"$earthly" --version
"$earthly" ./old-behaviour+myimage

docker inspect myimage:test >/dev/null
docker inspect mysubimage:test >/dev/null

echo "=== Test New Behaviour ==="

docker rm -f myimage:test
docker rm -f mysubimage:test
"$earthly" ./new-behaviour+myimage

docker inspect myimage:latest >/dev/null

if docker inspect mysubimage:latest; then
    echo "ERROR: mysubimage should not have been saved."
    exit 1
fi

echo "save-images test passed"
