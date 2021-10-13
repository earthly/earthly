#!/usr/bin/env bash
set -xeu

earthly=${earthly:=earthly}
earthly=$(realpath "$earthly")
echo "running tests with $earthly"

# prevent the self-update of earthly from running (this ensures no bogus data is printed to stdout,
# which would mess with the secrets data being fetched)
date +%s > /tmp/last-earthly-prerelease-check

# ensure earthly login works (and print out who gets logged in)
test -n "$EARTHLY_TOKEN"
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
    RUN echo "hello my world" > /data
    SAVE IMAGE earthly-export-test-3:test
EOF

"$earthly" prune --reset
"$earthly" +test3

docker run --rm earthly-export-test-3:test cat /data | grep "hello my world"


# Test 4: export multiplatform image
echo ==== Running test 4 ====
rm -rf /tmp/earthly-export-test-4
docker rmi earthly-export-test-4:test || true
docker rmi earthly-export-test-4:test_linux_amd64 || true
docker rmi earthly-export-test-4:test_linux_arm64 || true
docker rmi earthly-export-test-4:test_linux_arm_v7 || true

mkdir /tmp/earthly-export-test-4
cd /tmp/earthly-export-test-4
cat >> Earthfile <<EOF
multi4:
    BUILD --platform=linux/amd64 --platform=linux/arm64 --platform=linux/arm/v7 +test4

test4:
    FROM busybox:latest
    RUN echo "hello my world" > /data
    RUN uname -m >> /data
    SAVE IMAGE earthly-export-test-4:test
EOF

"$earthly" prune --reset
"$earthly" +multi4

docker run --rm earthly-export-test-4:test cat /data | grep "hello my world"
docker run --rm earthly-export-test-4:test cat /data | grep "$(uname -m)"
docker run --rm earthly-export-test-4:test_linux_amd64 cat /data | grep "hello my world"
docker run --rm earthly-export-test-4:test_linux_amd64 cat /data | grep "x86_64"
docker run --rm earthly-export-test-4:test_linux_arm64 cat /data | grep "hello my world"
docker run --rm earthly-export-test-4:test_linux_arm64 cat /data | grep "aarch64"
docker run --rm earthly-export-test-4:test_linux_arm_v7 cat /data | grep "hello my world"
docker run --rm earthly-export-test-4:test_linux_arm_v7 cat /data | grep "armv7l"


# Test 5: export multiple images
echo ==== Running test 5 ====
rm -rf /tmp/earthly-export-test-5
docker rmi earthly-export-test-5:test-img1 || true
docker rmi earthly-export-test-5:test-img2 || true

mkdir /tmp/earthly-export-test-5
cd /tmp/earthly-export-test-5
cat >> Earthfile <<EOF
all5:
    BUILD +test5-img1
    BUILD +test5-img2

test5-img1:
    FROM busybox:latest
    RUN echo "hello my world 1" > /data
    SAVE IMAGE earthly-export-test-5:test-img1

test5-img2:
    FROM busybox:latest
    RUN echo "hello my world 2" > /data
    SAVE IMAGE earthly-export-test-5:test-img2
EOF

"$earthly" prune --reset
"$earthly" +all5

docker run --rm earthly-export-test-5:test-img1 cat /data | grep "hello my world 1"
docker run --rm earthly-export-test-5:test-img2 cat /data | grep "hello my world 2"
