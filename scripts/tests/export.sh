#!/usr/bin/env bash
set -xeu

earthly=${earthly:=earthly}
if [ "$earthly" != "earthly" ]; then
  earthly=$(realpath "$earthly")
fi
echo "running tests with $earthly"
"$earthly" --version
frontend="${frontend:-$(which docker || which podman)}"
test -n "$frontend" || (>&2 echo "Error: frontend is empty" && exit 1)
echo "using frontend $frontend"

PATH="$(realpath "$(dirname "$0")/../acbtest"):$PATH"

# prevent the self-update of earthly from running (this ensures no bogus data is printed to stdout,
# which would mess with the secrets data being fetched)
date +%s > /tmp/last-earthly-prerelease-check

set +x # dont remove or the token will be leaked
if [ -z "${EARTHLY_TOKEN:-}" ]; then
  echo "using EARTHLY_TOKEN from earthly secrets"
  EARTHLY_TOKEN="$(earthly secrets --org earthly-technologies --project core get earthly-token-for-satellite-tests)"
  export EARTHLY_TOKEN
fi
test -n "$EARTHLY_TOKEN" || (echo "error: EARTHLY_TOKEN is not set" && exit 1)
set -x

EARTHLY_INSTALLATION_NAME="earthly-integration"
export EARTHLY_INSTALLATION_NAME
rm -rf "$HOME/.earthly.integration/"

echo "$earthly"
# ensure earthly login works (and print out who gets logged in)
"$earthly" account login

# Test 1: export without anything
echo ==== Running test 1 ====
rm -rf /tmp/earthly-export-test-1
"$frontend" rmi earthly-export-test-1:test || true

mkdir /tmp/earthly-export-test-1
cd /tmp/earthly-export-test-1
cat >> Earthfile <<EOF
VERSION 0.7
test1:
    FROM busybox:latest
    SAVE IMAGE earthly-export-test-1:test
EOF

"$earthly" prune --reset
"$earthly" +test1

"$frontend" run --rm earthly-export-test-1:test

# Test 2: export with only a CMD set
echo ==== Running test 2 ====
rm -rf /tmp/earthly-export-test-2
"$frontend" rmi earthly-export-test-2:test || true

mkdir /tmp/earthly-export-test-2
cd /tmp/earthly-export-test-2
cat >> Earthfile <<EOF
VERSION 0.7
test2:
    FROM busybox:latest
    CMD echo "running default cmd"
    SAVE IMAGE earthly-export-test-2:test
EOF

"$earthly" prune --reset
"$earthly" +test2

"$frontend" run --rm earthly-export-test-2:test | acbgrep "running default cmd"

# Test 3: export with a single RUN
echo ==== Running test 3 ====
rm -rf /tmp/earthly-export-test-3
"$frontend" rmi earthly-export-test-3:test || true

mkdir /tmp/earthly-export-test-3
cd /tmp/earthly-export-test-3
cat >> Earthfile <<EOF
VERSION 0.7
test3:
    FROM busybox:latest
    RUN echo "hello my world" > /data
    SAVE IMAGE earthly-export-test-3:test
EOF

"$earthly" prune --reset
"$earthly" +test3

"$frontend" run --rm earthly-export-test-3:test cat /data | acbgrep "hello my world"


# Test 4: export multiplatform image
echo ==== Running test 4 ====
rm -rf /tmp/earthly-export-test-4
"$frontend" rmi earthly-export-test-4:test || true
"$frontend" rmi earthly-export-test-4:test_linux_amd64 || true
"$frontend" rmi earthly-export-test-4:test_linux_arm64 || true
"$frontend" rmi earthly-export-test-4:test_linux_arm_v7 || true

mkdir /tmp/earthly-export-test-4
cd /tmp/earthly-export-test-4
cat >> Earthfile <<EOF
VERSION 0.7

multi4:
    # NOTE: keep amd64 in the middle, since earthly will fallback to the first defined platform
    # in case loadDockerManifest fails
    BUILD --platform=linux/arm/v7 --platform=linux/amd64 --platform=linux/arm64 +test4

test4:
    FROM busybox:latest
    RUN echo "hello my world" > /data
    RUN uname -m >> /data
    SAVE IMAGE earthly-export-test-4:test
EOF

"$earthly" prune --reset
"$earthly" +multi4

"$frontend" run --rm earthly-export-test-4:test cat /data | acbgrep "hello my world"
"$frontend" run --rm earthly-export-test-4:test cat /data | acbgrep "$(uname -m)"
"$frontend" run --rm earthly-export-test-4:test_linux_amd64 cat /data | acbgrep "hello my world"
"$frontend" run --rm earthly-export-test-4:test_linux_amd64 cat /data | acbgrep "x86_64"
"$frontend" run --rm earthly-export-test-4:test_linux_arm64 cat /data | acbgrep "hello my world"
"$frontend" run --rm earthly-export-test-4:test_linux_arm64 cat /data | acbgrep "aarch64"
"$frontend" run --rm earthly-export-test-4:test_linux_arm_v7 cat /data | acbgrep "hello my world"
"$frontend" run --rm earthly-export-test-4:test_linux_arm_v7 cat /data | acbgrep "armv7l"


# Test 5: export multiple images
echo ==== Running test 5 ====
rm -rf /tmp/earthly-export-test-5
"$frontend" rmi earthly-export-test-5:test-img1 || true
"$frontend" rmi earthly-export-test-5:test-img2 || true

mkdir /tmp/earthly-export-test-5
cd /tmp/earthly-export-test-5
cat >> Earthfile <<EOF
VERSION 0.7

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

"$frontend" run --rm earthly-export-test-5:test-img1 cat /data | acbgrep "hello my world 1"
"$frontend" run --rm earthly-export-test-5:test-img2 cat /data | acbgrep "hello my world 2"

# Test 6: no manifest list
echo ==== Running test 6 ====
rm -rf /tmp/earthly-export-test-6
"$frontend" rmi earthly-export-test-6:test || true
"$frontend" rmi earthly-export-test-6:test_linux_arm64 || true

mkdir /tmp/earthly-export-test-6
cd /tmp/earthly-export-test-6
cat >> Earthfile <<EOF
VERSION 0.7

multi6:
    BUILD --platform=linux/arm64 +test6

test6:
    FROM busybox:latest
    RUN echo "hello my world" > /data
    RUN uname -m >> /data
    SAVE IMAGE --no-manifest-list earthly-export-test-6:test
EOF

"$earthly" prune --reset
"$earthly" +multi6

"$frontend" run --rm earthly-export-test-6:test cat /data | acbgrep "hello my world"
"$frontend" run --rm earthly-export-test-6:test cat /data | acbgrep "aarch64"
if "$frontend" inspect earthly-export-test-6:test_linux_arm64 >/dev/null 2>&1 ; then
    echo "Expected failure"
    exit 1
fi

# Test 7: remote cache on target with only BUILDs
echo ==== Running test 7 ====
rm -rf /tmp/earthly-export-test-7
mkdir /tmp/earthly-export-test-7
cd /tmp/earthly-export-test-7
cat >> Earthfile <<EOF
VERSION 0.7
test7:
    BUILD +b
b:
    FROM busybox:latest
EOF

# This simply tests that this does not hang (#1945).
timeout -k 11m 10m "$earthly" --ci --push --remote-cache earthly/test-cache:export-test-7 +test7

# Test 8: Earthly LABELS
echo ==== Running test 8 ====
rm -rf /tmp/earthly-export-test-8
"$frontend" rmi earthly-export-test-8a:test || true
"$frontend" rmi earthly-export-test-8b:test || true

mkdir /tmp/earthly-export-test-8
cd /tmp/earthly-export-test-8
cat >> Earthfile <<EOF
VERSION 0.7

test8:
    FROM busybox:latest
    RUN echo "hello my world" > /data
    SAVE IMAGE earthly-export-test-8a:test
    LABEL foo=bar
    SAVE IMAGE earthly-export-test-8b:test
EOF

"$earthly" prune --reset
"$earthly" +test8

label_count=$("$frontend" inspect earthly-export-test-8a:test | jq .[].Config.Labels | grep -c dev.earthly.)
if [ "$label_count" -ne "3" ]; then
    echo "Expected 3 dev.earthly labels on first image; got $label_count"
    exit 1
fi

label_count=$("$frontend" inspect earthly-export-test-8b:test | jq .[].Config.Labels | grep -c dev.earthly.)
if [ "$label_count" -ne "3" ]; then
    echo "Expected 3 dev.earthly labels on second image; got $label_count"
    exit 1
fi

echo "=== All tests have passed ==="
