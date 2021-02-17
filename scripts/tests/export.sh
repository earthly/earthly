#!/usr/bin/env bash
set -eu

earthly=${earthly:=earthly}
earthly=$(realpath "$earthly")
echo "running tests with $earthly"

# prevent the self-update of earthly from running (this ensures no bogus data is printed to stdout,
# which would mess with the secrets data being fetched)
date +%s > /tmp/last-earthly-prerelease-check

# ensure earthly login works (and print out who gets logged in)
"$earthly" account login

# Test 1: export without anything
echo ==== Running test 1 ====
rm -rf /tmp/earthly-export-test-1
docker rmi earthly-export-test-1:test || true

mkdir /tmp/earthly-export-test-1
cd /tmp/earthly-export-test-1
cat >> Earthfile <<EOF
test1:
    FROM busybox:latest
    SAVE IMAGE earthly-export-test-1:test
EOF

"$earthly" prune --reset
"$earthly" +test1

docker run --rm earthly-export-test-1:test

# Test 2: export with only a CMD set
echo ==== Running test 2 ====
rm -rf /tmp/earthly-export-test-2
docker rmi earthly-export-test-2:test || true

mkdir /tmp/earthly-export-test-2
cd /tmp/earthly-export-test-2
cat >> Earthfile <<EOF
test2:
    FROM busybox:latest
    CMD echo "running default cmd"
    SAVE IMAGE earthly-export-test-2:test
EOF

"$earthly" prune --reset
"$earthly" +test2

docker run --rm earthly-export-test-2:test | grep "running default cmd"

# Test 3: export with a single RUN
echo ==== Running test 3 ====
rm -rf /tmp/earthly-export-test-3
docker rmi earthly-export-test-3:test || true

mkdir /tmp/earthly-export-test-3
cd /tmp/earthly-export-test-3
cat >> Earthfile <<EOF
test3:
    FROM busybox:latest
    RUN true # hack needed to make this work
    SAVE IMAGE earthly-export-test-3:test
EOF

"$earthly" prune --reset
"$earthly" +test3

docker run --rm earthly-export-test-3:test | grep "running default cmd"
