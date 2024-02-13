#!/bin/bash
# This test is designed to be run directly by github actions or on your host (i.e. not earthly-in-earthly)
set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
echo "using earthly=$(realpath "$earthly")"

rm .testdata || true # cleanup

"$earthly" $@ +test

test -f .testdata
test "$(cat .testdata)" = "pocus"
test "$(md5sum .testdata | awk '{print $1}')" = "18c27dcdffa2bfeca760684b1dfd72c8"
