#!/bin/bash
# This test is designed to be run directly by github actions or on your host (i.e. not earthly-in-earthly)
set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
echo "using earthly=$(realpath "$earthly")"

echo "running test part a"
rm .testbean || true # cleanup
rm .testpea || true # cleanup

"$earthly" $@ +test

test -f .testbean
test -f .testpea
test "$(cat .testbean)" = "garbanzo"
test "$(md5sum .testbean | awk '{print $1}')" = "cac81ad845e0333d7c1d73fbf36e5152"
test "$(cat .testpea)" = "chick"
test "$(md5sum .testpea | awk '{print $1}')" = "3278e0ac92da1d2d09e31d3698c85923"

echo "running test part b"
rm .testbean || true # cleanup
rm .testpea || true # cleanup

set +e
"$earthly" $@ +test --fail=yesplease
exit_code="$?"

test -f .testbean
test -f .testpea
test "$(cat .testbean)" = "garbanzo"
test "$(md5sum .testbean | awk '{print $1}')" = "cac81ad845e0333d7c1d73fbf36e5152"
test "$(cat .testpea)" = "chick"
test "$(md5sum .testpea | awk '{print $1}')" = "3278e0ac92da1d2d09e31d3698c85923"
